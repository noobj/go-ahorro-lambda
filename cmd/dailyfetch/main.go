package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/internal/middleware"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	EntryRepository "github.com/noobj/swim-crowd-lambda-go/internal/repositories/entry"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
)

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

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format(OutputFormat)
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local).Format(OutputFormat)

	startFromQuery, startExist := request.QueryStringParameters["start"]
	endFromQuery, endExist := request.QueryStringParameters["end"]

	if startExist && endExist {
		start = processTimeQueryString(startFromQuery, true)
		end = processTimeQueryString(endFromQuery, false)
	}

	var entryRepository repositories.Repository
	container.Resolve(&entryRepository)

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

	cursor := entryRepository.Aggregate([]bson.D{matchStage, groupStage})

	var results []EntryRepository.EntryGroup

	for cursor.Next(context.TODO()) {
		var result EntryRepository.EntryGroup
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
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
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
	container.Singleton(func() repositories.Repository {
		return &EntryRepository.EntryModel{}
	})

	lambda.Start(middleware.Logging(Handler))
}
