package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
)

type Specification struct {
	DynamoRandTable string `required:"true" split_words:"true"`
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](401)
	}

	env := config.GetInstance()

	config := helper.GenerateOauthConfig()

	token := oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  user.GoogleAccessToken,
		RefreshToken: user.GoogleRefreshToken,
		// TODO: use real value
		Expiry: time.Now().Add(time.Hour * -2),
	}

	client := config.Client(ctx, &token)
	service, _ := drive.New(client)

	randStateTable := env.DynamoRandTable

	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	authURL := config.AuthCodeURL(randState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(user.Id.Hex()),
			},
			"Randstate": {
				S: aws.String(randState),
			},
			"ttl": {
				N: aws.String(fmt.Sprintf("%d", time.Now().Add(time.Minute*5).Unix())),
			},
		},
		TableName: aws.String(randStateTable),
	}

	_, err := svc.PutItem(input)

	if err != nil {
		log.Printf("Dynamo insert randstate error: %v", err)
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](500)
	}

	file, err := service.Files.List().Q("name contains 'ahorro'").OrderBy("createdTime desc").PageSize(1).Do()

	if err != nil {
		log.Printf("Unable to create Drive service: %v", err)
		return helper.GenerateApiResponse[events.APIGatewayProxyResponse](authURL)
	}

	fileId := file.Files[0].Id

	if _, err = service.Files.Get(fileId).Do(); err != nil {
		log.Printf("Unable to create Drive service: %v", err)
		return helper.GenerateApiResponse[events.APIGatewayProxyResponse](authURL)
	}

	message := sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(user.Id.Hex()),
			},
		},
		MessageBody: aws.String("Sync ahorro entries with latest backup file"),
	}
	_, err = helper.SendSqsMessage(&message)
	if err != nil {
		log.Println("sending sqs error: ", err)
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](500)
	}

	dynamoTaskTable := env.DynamoTaskTable
	fmt.Printf("+%v", dynamoTaskTable)
	taskId := uuid.New()
	input = &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"TaskId": {
				S: aws.String(taskId.String()),
			},
			"Completed": {
				BOOL: aws.Bool(false),
			},
			"ttl": {
				N: aws.String(fmt.Sprintf("%d", time.Now().Add(time.Minute*10).Unix())),
			},
		},
		TableName: aws.String(dynamoTaskTable),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Printf("Dynamo Insert TaskId Error: %v", err)
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](500)
	}

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse](taskId)
}

func main() {
	lambda.Start(jwtMiddleWare.Auth(Handler))
}
