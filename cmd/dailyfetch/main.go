package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/noobj/swim-crowd-lambda-go/cmd/dailyfetch/models"
	"github.com/noobj/swim-crowd-lambda-go/internal/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const OutputFormat = "2006-01-02 15:04:05"

type Application struct {
	entryModel interface {
		FetchEntriesByTimeRange(start string, end string) []models.EntryGroup
	}
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

func (app *Application) Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format(OutputFormat)
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local).Format(OutputFormat)

	startFromQuery, startExist := request.QueryStringParameters["start"]
	endFromQuery, endExist := request.QueryStringParameters["end"]

	if startExist && endExist {
		start = processTimeQueryString(startFromQuery, true)
		end = processTimeQueryString(endFromQuery, false)
	}

	results := app.entryModel.FetchEntriesByTimeRange(start, end)

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
	app := &Application{entryModel: models.EntryModel{}}
	lambda.Start(middleware.Logging(app.Handler))
}
