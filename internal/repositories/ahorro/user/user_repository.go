package repository

import (
	"context"

	"github.com/noobj/go-serverless-services/internal/mongodb"
	"github.com/noobj/go-serverless-services/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id                 primitive.ObjectID `json:"_id" bson:"_id"`
	Account            string
	Password           string
	GoogleRefreshToken string `bson:"googleRefreshToken`
	GoogleAccessToken  string `bson:"googleAccessToken`
}

type UserRepository struct {
	repositories.IRepository
	Collection *mongo.Collection
}

//go:generate mockgen -source=user_repository.go -package repositories -aux_files repositories=../../repository.go -destination ../../mocks/user/mock_user_repository.go
type IUserRepository interface {
	repositories.IRepository
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

func New() *UserRepository {
	baseRepository := repositories.BaseRepository{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("ahorro").Collection("loginInfos"),
	}
	repo := UserRepository{IRepository: baseRepository}
	repo.Collection = baseRepository.Collection

	return &repo
}

func (u UserRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return u.Collection.UpdateOne(ctx, filter, update, opts...)
}
