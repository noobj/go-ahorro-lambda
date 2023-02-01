package repository

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Amount   float32            `json:"amount" bson:"amount"`
	Date     string             `json:"date"`
	Descr    string             `json:"descr"`
	Category primitive.ObjectID `json:"category,omitempty" bson:"category,omitempty"`
	User     primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
}

type Category struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string
	User  primitive.ObjectID
	Color string
	V     int `bson:"__v,omitempty"`
}

type EntryRepository struct {
	IRepository
}

func New() *EntryRepository {
	baseRepository := BaseRepository{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("ahorro").Collection("entries"),
	}
	repo := EntryRepository{IRepository: baseRepository}

	return &repo
}
