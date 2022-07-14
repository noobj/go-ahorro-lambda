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
}

func (m EntryModel) InsertOne(doc bson.D) {
	client := mongodb.GetInstance()

	coll := client.Database("swimCrowdDB").Collection("entries")
	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
}

func (m EntryModel) Aggregate(stages []bson.D) *mongo.Cursor {
	client := mongodb.GetInstance()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("swimCrowdDB").Collection("entries")

	cursor, err := coll.Aggregate(context.TODO(), stages)
	if err != nil {
		panic(err)
	}

	return cursor
}
