package jwt_auth

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	container "github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/middleware"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth[T middleware.ApiRequest](f middleware.HandlerFunc[T]) middleware.HandlerFunc[T] {
	return func(ctx context.Context, r T) (events.APIGatewayProxyResponse, error) {
		v2Request, ok := any(r).(events.APIGatewayV2HTTPRequest)
		if !ok {
			return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
		}
		cookiesMap := helper.ParseCookie(v2Request.Cookies)
		if _, ok := cookiesMap["access_token"]; !ok {
			return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
		}

		key := os.Getenv("ACCESS_TOKEN_SECRET")
		payload, err := helper.ExtractPayloadFromToken(key, cookiesMap["access_token"])
		if err != nil {
			return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
		}
		user, err := getUserForPayload(payload)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: "please login in", StatusCode: 401}, nil
		}

		ctxWithUser := context.WithValue(ctx, helper.ContextKeyUser, *user)

		return f(ctxWithUser, any(v2Request).(T))
	}
}

func getUserForPayload(payload interface{}) (*UserRepository.User, error) {
	userId, ok := payload.(string)
	userObjId, _ := primitive.ObjectIDFromHex(userId)
	if !ok {
		log.Printf("wrong payload format: %v", payload)
		return nil, fmt.Errorf("wrong payload format")
	}

	var userRepository repositories.IRepository
	container.NamedResolve(&userRepository, "UserRepo")
	var user UserRepository.User
	err := userRepository.FindOne(context.TODO(), bson.M{"_id": userObjId}).Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
