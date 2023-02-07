package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	bindioc "github.com/noobj/go-serverless-services/internal/middleware/bind-ioc"
	"github.com/noobj/go-serverless-services/internal/mongodb"
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

func (this Invoker) insertNewRefreshTokenIntoLoginInfo(userId primitive.ObjectID, refreshToken string) {
	loginInfo := LoginInfoRepository.LoginInfo{
		User:         userId,
		RefreshToken: refreshToken,
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}
	this.loginInfoRepository.InsertOne(context.TODO(), loginInfo)
}

type Invoker struct {
	userRepository      UserRepository.UserRepository           `container:"type"`
	loginInfoRepository LoginInfoRepository.LoginInfoRepository `container:"type"`
}

func (this Invoker) Invoke(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayV2HTTPResponse, error) {
	var requestBody LoginDto

	formData, err := helper.ParseMultipartForm(request.Headers["content-type"], strings.NewReader(request.Body), request.IsBase64Encoded)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{Body: "request body error", StatusCode: 400}, nil
	}

	requestBody.Account = formData.Value["account"][0]
	requestBody.Password = formData.Value["password"][0]
	if requestBody.Account == "" {
		return events.APIGatewayV2HTTPResponse{Body: "request body error", StatusCode: 400}, nil
	}

	var user UserRepository.User
	err = this.userRepository.FindOne(context.TODO(), bson.M{"account": requestBody.Account}).Decode(&user)
	if err != nil {
		log.Println(err)
		return events.APIGatewayV2HTTPResponse{Body: "couldn't find the user", StatusCode: 404}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		log.Println(err)
		return events.APIGatewayV2HTTPResponse{Body: "account and password not match", StatusCode: 404}, nil
	}

	token, err := helper.GenerateAccessToken(user.Id.Hex())
	if err != nil {
		log.Println("Couldn't generate access token", err)
		return events.APIGatewayV2HTTPResponse{Body: "internal error", StatusCode: 500}, nil
	}

	refreshToken, err := helper.GenerateRefreshToken(user.Id.Hex())
	if err != nil {
		log.Println("Couldn't generate refresh token", err)
		return events.APIGatewayV2HTTPResponse{Body: "internal error", StatusCode: 500}, nil
	}

	this.insertNewRefreshTokenIntoLoginInfo(user.Id, refreshToken)

	resp := events.APIGatewayV2HTTPResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "logged-in",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	env := config.GetInstance()
	accessTokenExpireTime := env.AccessTokenExpirationTime
	refreshTokenExpireTime := env.RefreshTokenExpirationTime
	cookieWithAccessTkn := http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(time.Second * time.Duration(accessTokenExpireTime)),
		Path:     "/",
	}
	cookieWithRefreshTkn := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(time.Second * time.Duration(refreshTokenExpireTime)),
		Path:     "/auth",
	}
	helper.SetCookie(cookieWithAccessTkn, &resp)
	helper.SetCookie(cookieWithRefreshTkn, &resp)

	return resp, nil
}

func main() {
	defer mongodb.Disconnect()()

	lambda.Start(bindioc.Handle[events.APIGatewayProxyRequest, events.APIGatewayV2HTTPResponse](&Invoker{}))
}
