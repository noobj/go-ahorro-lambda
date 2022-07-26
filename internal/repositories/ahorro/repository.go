package ahorro

import (
	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	. "github.com/noobj/swim-crowd-lambda-go/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Amount int
	Date   string
	Descr  string
}

type Category struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string
	User  primitive.ObjectID
	Color string
	V     int `bson:"__v"`
}

type AhorroRepository struct {
	AbstractRepository
}

func New() *AhorroRepository {
	abstractRepository := AbstractRepository{
		BaseRepository: BaseRepository{
			Client:     mongodb.GetInstance(),
			Collection: mongodb.GetInstance().Database("ahorro").Collection("entries"),
		},
	}
	repo := AhorroRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}
