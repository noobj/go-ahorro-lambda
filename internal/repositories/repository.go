package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=repository.go -destination mocks/mock_repository.go
type IRepository interface {
	Aggregate(interface{}) []bson.M
	// TODO: use generic for Entry type, not bson
	InsertOne(bson.D)
	Disconnect() func()
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
}

type BaseRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}
