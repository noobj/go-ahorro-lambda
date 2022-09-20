package jwt_auth

import (
	"context"
	"fmt"
	"log"

	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"github.com/noobj/jwtmiddleware"
	"github.com/noobj/jwtmiddleware/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth[T types.ApiRequest, R types.ApiResponse](f types.HandlerFunc[T, R]) types.HandlerFunc[T, R] {
	return jwtmiddleware.Handle(f, payloadHandler)
}

func payloadHandler(ctx context.Context, payload interface{}) (context.Context, error) {
	userId, ok := payload.(string)
	userObjId, _ := primitive.ObjectIDFromHex(userId)
	if !ok {
		log.Printf("wrong payload format: %v", payload)
		return nil, fmt.Errorf("wrong payload format")
	}

	userRepo := UserRepository.New()
	defer userRepo.Disconnect()()
	var user UserRepository.User

	err := userRepo.FindOne(context.TODO(), bson.M{"_id": userObjId}).Decode(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctx = context.WithValue(ctx, helper.ContextKeyUser, user)

	return ctx, nil
}
