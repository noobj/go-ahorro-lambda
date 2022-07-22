package entry

import (
	"context"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

type Entry struct {
	Amount int    `json:"amount"`
	Time   string `json:"time"`
}

type EntryRepository repositories.BaseRepository

func New() *EntryRepository {
	return &EntryRepository{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("swimCrowdDB").Collection("entries"),
	}
}

func (m EntryRepository) Disconnect() func() {
	return func() {
		if err := m.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func (m EntryRepository) InsertOne(doc bson.D) {
	_, err := m.Collection.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func (m EntryRepository) Aggregate(stages interface{}) []bson.M {
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
