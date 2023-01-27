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
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
)

var internalErrorhandler = func() (events.APIGatewayProxyResponse, error) {
	return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](500)
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	taskId, taskIdExist := request.QueryStringParameters["taskId"]
	if !taskIdExist {
		log.Println("Request has no task id")
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](400, "request query error")
	}

	env := config.GetInstance()

	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"TaskId": {
				S: aws.String(taskId),
			},
		},
		TableName: aws.String(env.DynamoTaskTable),
	}

	for i := 0; i < 20; i++ {
		item, err := svc.GetItem(input)

		if err != nil {
			fmt.Println("Fetch task status error", err)
			return internalErrorhandler()
		}

		if item.Item == nil {
			return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](400, "task id not found")
		}

		status := *item.Item["Completed"].N

		if status != fmt.Sprint(helper.Pending) {
			return helper.GenerateApiResponse[events.APIGatewayProxyResponse](status)
		}

		time.Sleep(time.Second)
	}

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse](helper.Pending)
}

func main() {
	lambda.Start(jwtMiddleWare.Handle(Handler))
}
