package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	bindioc "github.com/noobj/go-serverless-services/internal/middleware/bind-ioc"
	"github.com/noobj/go-serverless-services/internal/mongodb"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginDto struct {
	Account  string
	Password string
}

var errorHandler = helper.GenerateErrorResponse[events.APIGatewayV2HTTPResponse]

type Invoker struct {
	loginInfoRepository LoginInfoRepository.LoginInfoRepository `container:"type"`
}

func (this *Invoker) Invoke(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	cookiesMap := helper.ParseCookie(request.Cookies)

	if _, ok := cookiesMap["refresh_token"]; !ok {
		return errorHandler(401)
	}

	env := config.GetInstance()
	key := env.RefreshTokenSecret
	payload, err := helper.ExtractPayloadFromToken(key, cookiesMap["refresh_token"])
	if err != nil {
		return errorHandler(401)
	}
	userId, ok := payload.(string)
	if !ok {
		log.Printf("wrong payload format: %v", payload)
		return errorHandler(401)
	}

	var loginInfo LoginInfoRepository.LoginInfo
	this.loginInfoRepository.FindOne(context.TODO(), bson.M{"refreshToken": cookiesMap["refresh_token"]}).Decode(&loginInfo)

	if loginInfo.User.Hex() != userId {
		fmt.Println("Didn't match or find the loginInfo user", loginInfo.User.Hex(), userId)
		return errorHandler(401)
	}

	token, err := helper.GenerateAccessToken(userId)
	if err != nil {
		log.Println("Couldn't generate access token", err)
		return errorHandler(401)
	}

	resp := events.APIGatewayV2HTTPResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "refreshed",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	accessTokenExpireTime := env.AccessTokenExpirationTime
	cookieWithAccessTkn := http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Second * time.Duration(accessTokenExpireTime)),
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	helper.SetCookie(cookieWithAccessTkn, &resp)

	return resp, nil
}

func main() {
	defer mongodb.Disconnect()()

	lambda.Start(bindioc.Handle[events.APIGatewayV2HTTPRequest, events.APIGatewayV2HTTPResponse](&Invoker{}))
}
