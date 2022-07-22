package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	Aggregate(interface{}) []bson.M
	// TODO: use generic for Entry type, not bson
	InsertOne(bson.D)
	Disconnect() func()
}

type BaseRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}
