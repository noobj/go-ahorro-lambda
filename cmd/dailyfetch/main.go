package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/internal/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type EntryGroup struct {
	Date    string  `json:"date" bson:"_id,omitempty"`
	Entries []Entry `json:"entries" bson:"entries"`
}

type Entry struct {
	Amount int    `json:"amount"`
	Time   string `json:"time"`
}

const OutputFormat = "2006-01-02 15:04:05"

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

	results := fetchEntriesByTimeRange(start, end)

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
	lambda.Start(middleware.Logging(Handler))
}
