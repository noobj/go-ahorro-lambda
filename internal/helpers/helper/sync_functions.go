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

func PushSyncRequest(userId string) (events.APIGatewayProxyResponse, error) {
	env := config.GetInstance()
	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	taskId := uuid.New()

	message := sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(userId),
			},
			"TaskId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(taskId.String()),
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
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"TaskId": {
				S: aws.String(taskId.String()),
			},
			"Completed": {
				N: aws.String(fmt.Sprint(Pending)),
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

func UpdateTaskStatus(taskId string, status int) error {
	env := config.GetInstance()
	session, _ := session.NewSession()
	svc := dynamodb.New(session)

	item := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"TaskId": {
				S: aws.String(taskId),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {N: aws.String(fmt.Sprint(status))},
		},
		UpdateExpression: aws.String("SET Completed=:status"),
		TableName:        aws.String(env.DynamoTaskTable),
	}

	_, err := svc.UpdateItem(item)

	return err
}

const (
	Failed  = -1
	Pending = 0
	Done    = 1
)
