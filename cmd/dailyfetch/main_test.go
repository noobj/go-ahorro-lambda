package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/swim-crowd-lambda-go/cmd/dailyfetch/models"
)

type MockEntryModel struct{}

func (m MockEntryModel) FetchEntriesByTimeRange(start string, end string) []models.EntryGroup {
	return []models.EntryGroup{
		{
			Date: "2022-07-13",
			Entries: []models.Entry{
				{
					Amount: 1234,
					Time:   "2022-07-13 15:00",
				},
			},
		},
	}
}

func TestHandler(t *testing.T) {
	app := &Application{entryModel: MockEntryModel{}}
	res, err := app.Handler(context.TODO(), events.APIGatewayProxyRequest{})
	expectedRes := "[{\"date\":\"2022-07-13\",\"entries\":[{\"amount\":1234,\"time\":\"2022-07-13 15:00\"}]}]"

	if expectedRes != res.Body {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedRes, res.Body)
	}

	if err != nil {
		t.Errorf("error %s", err)
	}
}
