package ahorro

import (
	"context"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type EntryModel struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New() *EntryModel {
	return &EntryModel{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("ahorro").Collection("entries"),
	}
}

func (m EntryModel) Disconnect() func() {
	return func() {
		if err := m.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (m EntryModel) InsertOne(doc bson.D) {
	_, err := m.Collection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func (m EntryModel) Aggregate(stages interface{}) []bson.M {
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
