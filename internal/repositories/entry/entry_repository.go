package entry

import (
	"context"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EntryGroup struct {
	Date    string  `json:"date" bson:"_id,omitempty"`
	Entries []Entry `json:"entries" bson:"entries"`
}

type Entry struct {
	Amount int    `json:"amount"`
	Time   string `json:"time"`
}

type EntryModel struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New() *EntryModel {
	return &EntryModel{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("swimCrowdDB").Collection("entries"),
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

// TODO: should return EntryGroup[] or interface{}
// needed use generic
func (m EntryModel) Aggregate(stages []bson.D) []interface{} {
	cursor, err := m.Collection.Aggregate(context.TODO(), stages)
	if err != nil {
		panic(err)
	}

	var results []interface{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}
