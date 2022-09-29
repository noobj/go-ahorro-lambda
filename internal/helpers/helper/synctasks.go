package helper

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/noobj/go-serverless-services/internal/config"
)

func SyncTasks(userId string) (events.APIGatewayProxyResponse, error) {
	env := config.GetInstance()
	session, _ := session.NewSession()
	svc := dynamodb.New(session)

	message := sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(userId),
			},
		},
		MessageBody: aws.String("Sync ahorro entries with latest backup file"),
	}
	_, err := SendSqsMessage(&message)
	if err != nil {
		log.Println("sending sqs error: ", err)
		return GenerateErrorResponse[events.APIGatewayProxyResponse](500)
	}

	dynamoTaskTable := env.DynamoTaskTable
	fmt.Printf("+%v", dynamoTaskTable)
	taskId := uuid.New()
	input := &dynamodb.PutItemInput{
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
		return GenerateErrorResponse[events.APIGatewayProxyResponse](500)
	}

	return GenerateApiResponse[events.APIGatewayProxyResponse](taskId)
}
