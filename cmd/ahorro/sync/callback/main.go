package main

import (
	"context"
	"fmt"
	"log"

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
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var userRepositoryTmp repositories.IRepository
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
	}

	randState := struct {
		User  string
		State string
	}{}
	authRandStateCollection := mongodb.GetInstance().Database("ahorro").Collection("randState")
	err := authRandStateCollection.FindOne(ctx, bson.M{"user": user.Id.String()}).Decode(&randState)

	if err != nil {
		fmt.Println("fetch rand state error", err)
		return helper.GenerateInternalErrorResponse()
	}

	if randState.State != request.QueryStringParameters["state"] {
		return events.APIGatewayProxyResponse{Body: "rand state error", StatusCode: 401}, nil
	}

	config := helper.GenerateOauthConfig()

	token, err := config.Exchange(ctx, request.QueryStringParameters["code"], oauth2.AccessTypeOffline)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "exchange code error", StatusCode: 401}, nil
	}

	container.NamedResolve(&userRepositoryTmp, "UserRepo")
	userRepository, ok := userRepositoryTmp.(UserRepository.IUserRepository)

	if !ok {
		log.Println("resolve repository error")
		return helper.GenerateInternalErrorResponse()
	}

	_, err = userRepository.UpdateOne(context.TODO(), bson.M{"account": user.Account}, bson.M{"$set": bson.M{"googleAccessToken": token.AccessToken, "googleRefreshToken": token.RefreshToken}})
	if err != nil {
		log.Println("update error", err)
		return helper.GenerateInternalErrorResponse()
	}

	message := sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(user.Id.String()),
			},
		},
		MessageBody: aws.String("Sync ahorro entries with latest backup file"),
	}
	_, err = helper.SendSqsMessage(&message)
	if err != nil {
		log.Println("sending sqs error: ", err)
		return helper.GenerateInternalErrorResponse()
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
