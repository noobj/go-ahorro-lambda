package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	Aggregate([]bson.D) []interface{}
	InsertOne(bson.D)
	Disconnect() func()
}
