package ahorro

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Amount   string             `json:"amount" bson:"amount"`
	Date     string             `json:"date"`
	Descr    string             `json:"descr"`
	Category primitive.ObjectID `json:"category" bson:"category,omitempty"`
	User     primitive.ObjectID `json:"user" bson:"user,omitempty"`
}

type Category struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string
	User  primitive.ObjectID
	Color string
	V     int `bson:"__v,omitempty"`
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
