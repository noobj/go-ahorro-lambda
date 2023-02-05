package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	bindioc "github.com/noobj/go-serverless-services/internal/middleware/bind-ioc"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"github.com/noobj/go-serverless-services/internal/mongodb"
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

type Invoker struct {
	userRepository UserRepository.UserRepository `container:"type"`
}

func (this *Invoker) Invoke(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return authErrorhandler()
	}

	env := config.GetInstance()

	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(user.Id.Hex()),
			},
		},
		TableName: aws.String(env.DynamoRandTable),
	}

	item, err := svc.GetItem(input)

	if err != nil {
		fmt.Println("fetch rand state error", err)
		return internalErrorhandler()
	}

	if *item.Item["Randstate"].S != request.QueryStringParameters["state"] {
		return authErrorhandler("rand state error")
	}

	config := helper.GenerateOauthConfig()

	token, err := config.Exchange(ctx, request.QueryStringParameters["code"], oauth2.AccessTypeOffline)

	if err != nil {
		fmt.Println(err)
		return authErrorhandler("exchange code error")
	}

	if !ok {
		log.Println("resolve repository error")
		return internalErrorhandler()
	}

	_, err = this.userRepository.UpdateOne(context.TODO(), bson.M{"account": user.Account}, bson.M{"$set": bson.M{"googleAccessToken": token.AccessToken, "googleRefreshToken": token.RefreshToken}})
	if err != nil {
		log.Println("update error", err)
		return internalErrorhandler()
	}

	res, _ := helper.PushSyncRequest(user.Id.Hex())

	if res.StatusCode != 200 {
		return res, nil
	}

	return helper.GenerateRedirectResponse[events.APIGatewayProxyResponse](env.FrontendUrl)
}

func main() {
	defer mongodb.Disconnect()()
	invoker := Invoker{}

	lambda.Start(jwtMiddleWare.Handle(bindioc.Handle(invoker.Invoke, &invoker)))
}
