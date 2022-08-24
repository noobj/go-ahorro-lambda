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
	repositories.AbstractRepository
}

//go:generate mockgen -source=user_repository.go -package repositories -aux_files repositories=../../repository.go -destination mocks/mock_user_repository.go
type IUserRepository interface {
	repositories.IRepository
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

func New() *UserRepository {
	abstractRepository := repositories.AbstractRepository{
		BaseRepository: repositories.BaseRepository{
			Client:     mongodb.GetInstance(),
			Collection: mongodb.GetInstance().Database("ahorro").Collection("users"),
		},
	}
	repo := UserRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}

func (u UserRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return u.Collection.UpdateOne(ctx, filter, update, opts...)
}
