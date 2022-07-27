package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/repositories"
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

	accessTokenExpireTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION_TIME"))
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	if err != nil {
		log.Println("ENV ACCESS_TOKEN_EXPIRATION_TIME format wrong")
		return helper.GenerateInternalErrorResponse()
	}

	token, err := helper.GenerateJwtToken(user.Id.Hex(), accessTokenExpireTime, accessTokenSecret)
	if err != nil {
		log.Println("Couldn't generate access token", err)
		return helper.GenerateInternalErrorResponse()
	}
	// TODO: generate refresh tokens
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
