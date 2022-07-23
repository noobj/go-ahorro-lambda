package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	container "github.com/golobby/container/v3"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockEntryModel struct{}

func (m MockEntryModel) InsertOne(doc bson.D) {}
func (m MockEntryModel) Disconnect() func() {
	return func() {}
}

func (m MockEntryModel) Aggregate(stages interface{}) []bson.M {
	fakeObjId, _ := primitive.ObjectIDFromHex("62badc82d420270009a51019")

	fakeData := []bson.M{
		{
			"sum": 110,
			"_id": fakeObjId,
			"category": []bson.M{
				{
					"_id":   fakeObjId,
					"color": "#a4e56c",
					"name":  "Food",
					"user":  fakeObjId,
				},
			},
			"entries": []bson.M{
				{
					"_id":    fakeObjId,
					"amount": 110,
					"date":   "2022-01-05",
					"descr":  "fuck",
				},
			},
		},
		{
			"sum": 90,
			"_id": fakeObjId,
			"category": []bson.M{
				{
					"_id":   fakeObjId,
					"color": "#a4e51c",
					"name":  "Abc",
					"user":  fakeObjId,
				},
			},
			"entries": []bson.M{
				{
					"_id":    fakeObjId,
					"amount": 90,
					"date":   "2022-01-05",
					"descr":  "fuck",
				},
			},
		},
	}

	return fakeData
}

func TestHandlerPass(t *testing.T) {
	container.Singleton(func() repositories.IRepository {
		return MockEntryModel{}
	})

	var fakeRequest events.APIGatewayProxyRequest
	fakeRequest.QueryStringParameters = make(map[string]string)
	fakeRequest.QueryStringParameters["timeStart"] = "2022-01-01"
	fakeRequest.QueryStringParameters["timeEnd"] = "2022-01-31"
	res, err := Handler(context.TODO(), fakeRequest)

	expectedRes := "{\"categories\":[{\"_id\":\"62badc82d420270009a51019\",\"sum\":110,\"percentage\":\"0.55\",\"name\":\"Food\",\"entries\":[{\"_id\":\"62badc82d420270009a51019\",\"Amount\":110,\"Date\":\"2022-01-05\",\"Descr\":\"fuck\"}],\"color\":\"#a4e56c\"},{\"_id\":\"62badc82d420270009a51019\",\"sum\":90,\"percentage\":\"0.45\",\"name\":\"Abc\",\"entries\":[{\"_id\":\"62badc82d420270009a51019\",\"Amount\":90,\"Date\":\"2022-01-05\",\"Descr\":\"fuck\"}],\"color\":\"#a4e51c\"}],\"total\":200}"

	if expectedRes != res.Body {
		t.Errorf("\n...expected = %v\n...obtained = %v", expectedRes, res.Body)
	}

	if err != nil {
		t.Errorf("error %s", err)
	}
}

func TestPanicWithNoQueryStringParam(t *testing.T) {
	expectedRes := "something wrong with time query string"

	defer func() {
		if r := recover(); r != nil {
			if r != expectedRes {
				t.Errorf("\n...expected = %v\n...obtained = %v", expectedRes, r)
			}
		}
	}()

	container.Singleton(func() repositories.IRepository {
		return MockEntryModel{}
	})

	Handler(context.TODO(), events.APIGatewayProxyRequest{})
}
