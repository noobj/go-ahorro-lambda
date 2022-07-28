package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/repositories"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type LoginDto struct {
	Account  string
	Password string
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var ahorroRepository repositories.IRepository
	var requestBody LoginDto

	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "request body error", StatusCode: 404}, nil
	}

	container.Resolve(&ahorroRepository)
	defer ahorroRepository.Disconnect()()

	var user UserRepository.User
	err = ahorroRepository.FindOne(context.TODO(), bson.M{"account": requestBody.Account}).Decode(&user)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Couldn't find the user", StatusCode: 404}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "account and password not match", StatusCode: 404}, nil
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return helper.GenerateInternalErrorResponse()
	}

	token, err := helper.GenerateAccessToken(user.Id.Hex())
	if err != nil {
		log.Println("Couldn't generate access token", err)
		return helper.GenerateInternalErrorResponse()
	}

	refreshToken, err := helper.GenerateRefreshToken(user.Id.Hex())
	if err != nil {
		log.Println("Couldn't generate refresh token", err)
		return helper.GenerateInternalErrorResponse()
	}

	fmt.Println(refreshToken)

	loginInfo := LoginInfoRepository.LoginInfo{
		User:         user.Id,
		RefreshToken: refreshToken,
	}
	loginInfoRepository := LoginInfoRepository.New()
	loginInfoRepository.InsertOne(loginInfo)

	// TODO: insert refresh token into LoginInfo
	// TODO: set up response cookies

	return helper.GenerateApiResponse(token)
}

func main() {
	container.Singleton(func() repositories.IRepository {
		return UserRepository.New()
	})

	lambda.Start(Handler)
}
