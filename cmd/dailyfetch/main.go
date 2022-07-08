package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response events.APIGatewayProxyResponse

type EntryGroup struct {
	Date    string  `json:"date" bson:"_id,omitempty"`
	Entries []Entry `json:"entries" bson:"entries"`
}

type Entry struct {
	Amount int    `json:"amount"`
	Time   string `json:"time"`
}

const OutputFormat = "2006-01-02 15:04:05"

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

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	client := mongodb.InitMongoDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format(OutputFormat)
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local).Format(OutputFormat)

	startFromQuery, startExist := request.QueryStringParameters["start"]
	endFromQuery, endExist := request.QueryStringParameters["end"]

	if startExist && endExist {
		start = processTimeQueryString(startFromQuery, true)
		end = processTimeQueryString(endFromQuery, false)
	}

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

	var buf bytes.Buffer

	body, err := json.Marshal(results)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
