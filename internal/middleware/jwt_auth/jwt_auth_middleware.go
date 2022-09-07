package jwt_auth

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/middleware"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth[T middleware.ApiRequest, R helper.ApiResponse](f middleware.HandlerFunc[T, R]) middleware.HandlerFunc[T, R] {
	return func(ctx context.Context, r T) (R, error) {
		v2Request, ok := any(r).(events.APIGatewayV2HTTPRequest)
		if !ok {
			return helper.GenerateErrorResponse[R](401)
		}
		cookiesMap := helper.ParseCookie(v2Request.Cookies)
		if _, ok := cookiesMap["access_token"]; !ok {
			return helper.GenerateErrorResponse[R](401)
		}

		key := os.Getenv("ACCESS_TOKEN_SECRET")
		payload, err := helper.ExtractPayloadFromToken(key, cookiesMap["access_token"])
		if err != nil {
			return helper.GenerateErrorResponse[R](401)
		}
		user, err := getUserForPayload(payload)
		if err != nil {
			return helper.GenerateErrorResponse[R](401)
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
