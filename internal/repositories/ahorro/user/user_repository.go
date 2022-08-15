package repository

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Account  string
	Password string
}

type UserRepository struct {
	AbstractRepository
}

func New() *UserRepository {
	abstractRepository := AbstractRepository{
		BaseRepository: BaseRepository{
			Client:     mongodb.GetInstance(),
			Collection: mongodb.GetInstance().Database("ahorro").Collection("users"),
		},
	}
	repo := UserRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}
