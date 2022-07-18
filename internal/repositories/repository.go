package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	Aggregate([]bson.D) []any
	// TODO: use generic for Entry type, not bson
	InsertOne(bson.D)
	Disconnect() func()
}
