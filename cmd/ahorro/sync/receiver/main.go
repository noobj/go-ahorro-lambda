package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func sendSqsMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	qURL := os.Getenv("SQS_URL")
	input.QueueUrl = &qURL

	result, err := svc.SendMessage(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
	}
	fmt.Println(user)

	config := &oauth2.Config{
		ClientID:     "719667836696-c17u3p1dhu6832dp12ovdjk6tfr5qquc.apps.googleusercontent.com",
		ClientSecret: "92DOsm_TQr3B_xS7647hG6N2",
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/drive.readonly"},
	}

	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	config.RedirectURL = "https://ahorrojs.io/sync/callback"
	authURL := config.AuthCodeURL(randState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println(authURL)

	message := sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": {
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
	}
	result, err := sendSqsMessage(&message)
	if err != nil {
		log.Println("sending sqs error: ", err)
		return helper.GenerateInternalErrorResponse()
	}

	return helper.GenerateApiResponse(result)
}

func main() {
	userRepo := UserRepository.New()
	defer userRepo.Disconnect()()
	container.NamedSingleton("UserRepo", func() repositories.IRepository {
		return userRepo
	})

	lambda.Start(jwtMiddleWare.Auth(Handler))
}
