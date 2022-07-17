package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	// TODO: use generic
	Aggregate([]bson.D) []any
	InsertOne(bson.D)
	Disconnect() func()
}
