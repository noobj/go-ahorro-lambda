package ahorro

import (
	"context"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"

	"go.mongodb.org/mongo-driver/bson"
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

type AhorroRepository repositories.BaseRepository

func New() *AhorroRepository {
	return &AhorroRepository{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("ahorro").Collection("entries"),
	}
}

func (m AhorroRepository) Disconnect() func() {
	return func() {
		if err := m.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (m AhorroRepository) InsertOne(doc bson.D) {
	_, err := m.Collection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func (m AhorroRepository) Aggregate(stages interface{}) []bson.M {
	cursor, err := m.Collection.Aggregate(context.TODO(), stages)
	if err != nil {
		panic(err)
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}
