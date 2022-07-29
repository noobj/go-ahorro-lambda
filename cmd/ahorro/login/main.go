package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/repositories"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type LoginDto struct {
	Account  string
	Password string
}

func insertNewRefreshTokenIntoLoginInfo(userId primitive.ObjectID, refreshToken string) {
	loginInfo := LoginInfoRepository.LoginInfo{
		User:         userId,
		RefreshToken: refreshToken,
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}
	var loginInfoRepository repositories.IRepository
	container.NamedResolve(&loginInfoRepository, "LoginInfoRepo")
	loginInfoRepository.InsertOne(loginInfo)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userRepository repositories.IRepository
	var requestBody LoginDto

	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "request body error", StatusCode: 404}, nil
	}

	container.NamedResolve(&userRepository, "UserRepo")

	var user UserRepository.User
	err = userRepository.FindOne(context.TODO(), bson.M{"account": requestBody.Account}).Decode(&user)
	if err != nil {
		log.Panicln(err)
		return events.APIGatewayProxyResponse{Body: "Couldn't find the user", StatusCode: 404}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		log.Panicln(err)
		return events.APIGatewayProxyResponse{Body: "account and password not match", StatusCode: 404}, nil
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
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

	insertNewRefreshTokenIntoLoginInfo(user.Id, refreshToken)

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
	refreshTokenExpireTime, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_TIME"))
	cookieWithAccessTkn := http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Second * time.Duration(accessTokenExpireTime)),
	}
	cookieWithRefreshTkn := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Second * time.Duration(refreshTokenExpireTime)),
		Path:     "/auth",
	}
	helper.SetCookie(cookieWithAccessTkn, &resp)
	helper.SetCookie(cookieWithRefreshTkn, &resp)

	return resp, nil
}

func main() {
	userRepo := UserRepository.New()
	container.NamedSingleton("UserRepo", func() repositories.IRepository {
		return userRepo
	})

	container.NamedSingleton("LoginInfoRepo", func() repositories.IRepository {
		return LoginInfoRepository.New()
	})
	defer userRepo.Disconnect()()

	lambda.Start(Handler)
}
