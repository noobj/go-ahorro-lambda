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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response events.APIGatewayProxyResponse

type Entry struct {
	Id     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount int
	Time   string
}

func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	client := mongodb.InitMongoDB()
	fmt.Println(request.QueryStringParameters["startDate"])
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	timeFormat := "2006-01-02 15:4:5"
	t := time.Now()
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format(timeFormat)
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local).Format(timeFormat)
	coll := client.Database("swimCrowdDB").Collection("entries")
	matchStage := bson.D{{"$match", bson.D{
		{"$and",
			bson.A{
				bson.D{{"time", bson.D{{"$gt", start}}}},
				bson.D{{"time", bson.D{{"$lte", end}}}},
			},
		},
	}}}
	fmt.Println(matchStage)

	filter := bson.D{{"time", bson.D{{"$gt", "2022-07-01"}}}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []Entry
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"data": results,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
	// fmt.Println(Handler())
}
