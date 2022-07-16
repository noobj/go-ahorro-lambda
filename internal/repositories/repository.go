package repositories

import (
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories/entry"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	// TODO: use generic
	Aggregate([]bson.D) []entry.EntryGroup
	InsertOne(bson.D)
	Disconnect() func()
}
