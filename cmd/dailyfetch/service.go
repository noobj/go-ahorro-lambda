package main

import (
	"context"
	"fmt"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func processTimeQueryString(tString string, start bool) string {
	timeFormat := "2006-01-02"
	parsedTime, err := time.Parse(timeFormat, tString)
	if err != nil {
		panic(fmt.Sprintf("Could not parse time\n %s", err))
	}

	if start {
		return time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, time.Local).Format(OutputFormat)
	} else {
		return time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 23, 59, 59, 999999999, time.Local).Format(OutputFormat)
	}
}

func fetchEntriesByTimeRange(start string, end string) []EntryGroup {
	client := mongodb.ClientInstance
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
