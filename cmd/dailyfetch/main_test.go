package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	container "github.com/golobby/container/v3"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	EntryRepositories "github.com/noobj/swim-crowd-lambda-go/internal/repositories/entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockEntryModel struct{}

func (m MockEntryModel) InsertOne(doc bson.D) {}
func (m MockEntryModel) Disconnect() func() {
	return func() {}
}

func (m MockEntryModel) Aggregate(stages []bson.D) *mongo.Cursor {
	fakeData := []EntryRepositories.EntryGroup{
		{
			Date: "2022-07-13",
			Entries: []EntryRepositories.Entry{
				{
					Amount: 1234,
					Time:   "2022-07-13 15:00",
				},
			},
		},
	}

	mockData := make([]interface{}, len(fakeData))

	for i := range fakeData {
		mockData[i] = fakeData[i]
	}

	cur, _ := mongo.NewCursorFromDocuments(mockData, nil, nil)
	return cur
}

func TestHandler(t *testing.T) {
	container.Singleton(func() repositories.Repository {
		return MockEntryModel{}
	})
	res, err := Handler(context.TODO(), events.APIGatewayProxyRequest{})

	expectedRes := "[{\"date\":\"2022-07-13 (Wed)\",\"entries\":[{\"amount\":1234,\"time\":\"2022-07-13 15:00\"}]}]"

	if expectedRes != res.Body {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedRes, res.Body)
	}

	if err != nil {
		t.Errorf("error %s", err)
	}
}
