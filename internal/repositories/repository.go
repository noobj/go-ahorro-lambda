package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	db         string
	collection string
}

type Repository interface {
	Aggregate([]bson.D) *mongo.Cursor
	InsertOne(bson.D)
}
