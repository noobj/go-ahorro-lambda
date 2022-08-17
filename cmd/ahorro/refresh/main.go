package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/repositories"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginDto struct {
	Account  string
	Password string
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	cookiesMap := helper.ParseCookie(request.Cookies)

	if _, ok := cookiesMap["refresh_token"]; !ok {
		return helper.GenerateAuthErrorResponse()
	}
	key := os.Getenv("REFRESH_TOKEN_SECRET")
	payload, err := helper.ExtractPayloadFromToken(key, cookiesMap["refresh_token"])
	if err != nil {
		return helper.GenerateAuthErrorResponse()
	}
	userId, ok := payload.(string)
	if !ok {
		log.Printf("wrong payload format: %v", payload)
		return helper.GenerateAuthErrorResponse()
	}

	var loginInfoRepository repositories.IRepository
	var loginInfo LoginInfoRepository.LoginInfo
	container.NamedResolve(&loginInfoRepository, "LoginInfoRepo")
	loginInfoRepository.FindOne(context.TODO(), bson.M{"refreshToken": cookiesMap["refresh_token"]}).Decode(&loginInfo)

	if loginInfo.User.Hex() != userId {
		fmt.Println(loginInfo.User.Hex(), userId)
		return helper.GenerateAuthErrorResponse()
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
	}

	token, err := helper.GenerateAccessToken(userId)
	if err != nil {
		log.Println("Couldn't generate access token", err)
		return helper.GenerateInternalErrorResponse()
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		MultiValueHeaders: map[string][]string{
			"set-cookie": nil,
		},
	}

	accessTokenExpireTime, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION_TIME"))
	cookieWithAccessTkn := http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Second * time.Duration(accessTokenExpireTime)),
		Path:     "/",
	}
	helper.SetCookie(cookieWithAccessTkn, &resp)

	return resp, nil
}

func main() {
	repo := LoginInfoRepository.New()

	container.NamedSingleton("LoginInfoRepo", func() repositories.IRepository {
		return repo
	})
	defer repo.Disconnect()()

	lambda.Start(Handler)
}
