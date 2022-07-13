package models

import (
	"context"
	"time"

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

type EntryModel struct{}

func (m EntryModel) FetchEntriesByTimeRange(start string, end string) []EntryGroup {
	client := mongodb.GetInstance()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("swimCrowdDB").Collection("entries")
	matchStage := bson.D{{"$match", bson.D{
		{"$and",
			bson.A{
				bson.D{{"time", bson.D{{"$gt", start}}}},
				bson.D{{"time", bson.D{{"$lte", end}}}},
			},
		},
	}}}

	groupStage := bson.D{{
		"$group", bson.D{
			{
				"_id", bson.D{{
					"$substr", bson.A{"$time", 0, 10},
				}},
			},
			{
				"entries", bson.D{{
					"$push", bson.D{
						{"amount", "$amount"},
						{"time", "$time"},
					},
				}},
			},
		},
	}}

	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var results []EntryGroup

	for cursor.Next(context.TODO()) {
		var result EntryGroup
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}

		parsedDate, err := time.Parse("2006-01-02", result.Date)
		if err != nil {
			panic(err)
		}

		result.Date = parsedDate.Format("2006-01-02 (Mon)")
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

	return results
}
