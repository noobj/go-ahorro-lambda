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

var authErrorhandler = func(message ...string) (events.APIGatewayProxyResponse, error) {
	return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](401, message...)
}

var internalErrorhandler = func() (events.APIGatewayProxyResponse, error) {
	return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](500)
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var userRepositoryTmp repositories.IRepository
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return authErrorhandler()
	}

	randState := struct {
		User  string
		State string
	}{}
	authRandStateCollection := mongodb.GetInstance().Database("ahorro").Collection("randState")
	err := authRandStateCollection.FindOne(ctx, bson.M{"user": user.Id.String()}).Decode(&randState)

	if err != nil {
		fmt.Println("fetch rand state error", err)
		return internalErrorhandler()
	}

	if randState.State != request.QueryStringParameters["state"] {
		return authErrorhandler("rand state error")
	}

	config := helper.GenerateOauthConfig()

	token, err := config.Exchange(ctx, request.QueryStringParameters["code"], oauth2.AccessTypeOffline)

	if err != nil {
		fmt.Println(err)
		return authErrorhandler("exchange code error")
	}

	container.NamedResolve(&userRepositoryTmp, "UserRepo")
	userRepository, ok := userRepositoryTmp.(UserRepository.IUserRepository)

	if !ok {
		log.Println("resolve repository error")
		return internalErrorhandler()
	}

	_, err = userRepository.UpdateOne(context.TODO(), bson.M{"account": user.Account}, bson.M{"$set": bson.M{"googleAccessToken": token.AccessToken, "googleRefreshToken": token.RefreshToken}})
	if err != nil {
		log.Println("update error", err)
		return internalErrorhandler()
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
		return internalErrorhandler()
	}

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("ok")
}

func main() {
	userRepo := UserRepository.New()
	defer userRepo.Disconnect()()
	container.NamedSingleton("UserRepo", func() repositories.IRepository {
		return userRepo
	})

	lambda.Start(jwtMiddleWare.Auth(Handler))
}
