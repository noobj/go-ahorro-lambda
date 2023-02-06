package repository

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	"github.com/noobj/go-serverless-services/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string
	User  primitive.ObjectID
	Color string
}

type CategoryRepository struct {
	repositories.IRepository
}

func New() *CategoryRepository {
	baseRepository := repositories.BaseRepository{
		Collection: mongodb.GetInstance().Database("ahorro").Collection("categories"),
	}
	repo := CategoryRepository{IRepository: baseRepository}

	return &repo
}
