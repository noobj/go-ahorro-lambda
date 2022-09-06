package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"github.com/noobj/go-serverless-services/internal/mongodb"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	drive "google.golang.org/api/drive/v3"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
	}

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

	file, err := service.Files.List().Q("name contains 'ahorro'").OrderBy("createdTime desc").PageSize(1).Do()

	if err != nil {
		log.Printf("Unable to create Drive service: %v", err)
		randState := fmt.Sprintf("st%d", time.Now().UnixNano())
		authURL := config.AuthCodeURL(randState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

		authRandStateCollection := mongodb.GetInstance().Database("ahorro").Collection("randState")
		authRandStateCollection.InsertOne(ctx, bson.M{
			"user":  user.Id.String(),
			"state": randState,
		})
		return helper.GenerateApiResponse(authURL)
	}

	fileId := file.Files[0].Id

	if _, err = service.Files.Get(fileId).Do(); err != nil {
		log.Printf("Unable to create Drive service: %v", err)
		randState := fmt.Sprintf("st%d", time.Now().UnixNano())
		authURL := config.AuthCodeURL(randState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		return helper.GenerateApiResponse(authURL)
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
		return helper.GenerateInternalErrorResponse[events.APIGatewayProxyResponse]()
	}

	return helper.GenerateApiResponse("ok")
}

func main() {
	userRepo := UserRepository.New()
	defer userRepo.Disconnect()()
	container.NamedSingleton("UserRepo", func() repositories.IRepository {
		return userRepo
	})

	lambda.Start(jwtMiddleWare.Auth(Handler))
}
