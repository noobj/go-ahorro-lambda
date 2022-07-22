package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	container "github.com/golobby/container/v3"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

type MockEntryModel struct{}

func (m MockEntryModel) InsertOne(doc bson.D) {}
func (m MockEntryModel) Disconnect() func() {
	return func() {}
}

func (m MockEntryModel) Aggregate(stages interface{}) []bson.M {
	fakeData := []bson.M{
		{"Date": "2022-07-13",
			"Entries": []bson.M{
				{
					"Amount": 1234,
					"Time":   "2022-07-13 15:00",
				},
			},
		},
	}

	// var results []any

	// for _, data := range fakeData {
	// 	results = append(results, data)
	// }

	return fakeData
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
