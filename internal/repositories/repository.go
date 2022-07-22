package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	Aggregate(interface{}) []bson.M
	// TODO: use generic for Entry type, not bson
	InsertOne(bson.D)
	Disconnect() func()
}
