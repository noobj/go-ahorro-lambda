package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtMiddleWare "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found", err)
	}
	var userRepositoryTmp repositories.IRepository
	user, ok := helper.GetUserFromContext(ctx)
	if !ok {
		return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	config := &oauth2.Config{
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveReadonlyScope},
		RedirectURL:  "https://ahorrojs.io/sync/callback",
	}

	token, err := config.Exchange(ctx, request.QueryStringParameters["code"], oauth2.AccessTypeOffline)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{Body: "internal error", StatusCode: 500}, nil
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
