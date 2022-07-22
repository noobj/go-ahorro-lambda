package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/internal/middleware"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	EntryRepository "github.com/noobj/swim-crowd-lambda-go/internal/repositories/swim/entry"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
)

const OutputFormat = "2006-01-02 15:04:05"

type EntryGroup struct {
	Date    string                  `json:"date"`
	Entries []EntryRepository.Entry `json:"entries"`
}

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

	var entryRepository repositories.IRepository
	container.Resolve(&entryRepository)
	defer entryRepository.Disconnect()()

	matchStage := bson.M{"$match": bson.M{
		"$and": bson.A{
			bson.M{"time": bson.M{"$gt": start}},
			bson.M{"time": bson.M{"$lte": end}},
		},
	}}

	groupStage := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"$substr": bson.A{"$time", 0, 10},
			},
			"entries": bson.M{
				"$push": bson.M{
					"amount": "$amount",
					"time":   "$time",
				},
			},
		},
	}

	repoResults := entryRepository.Aggregate([]bson.M{matchStage, groupStage})
	var results []EntryGroup

	for _, repoResult := range repoResults {
		var result EntryGroup
		doc, _ := bson.Marshal(repoResult)
		err := bson.Unmarshal(doc, &result)
		if err != nil {
			panic("Repository returns wrong type")
		}
		parsedDate, err := time.Parse("2006-01-02", result.Date)
		if err != nil {
			panic(err)
		}

		result.Date = parsedDate.Format("2006-01-02 (Mon)")
		results = append(results, result)
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
	container.Singleton(func() repositories.IRepository {
		return EntryRepository.New()
	})

	lambda.Start(middleware.Logging(Handler))
}
