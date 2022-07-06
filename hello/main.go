package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Response events.APIGatewayProxyResponse

func initMongoDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	return client
}

func Handler(ctx context.Context) (Response, error) {
	client := initMongoDB()
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
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"time", bson.D{{"$gt", start}}}},
				bson.D{{"time", bson.D{{"$lte", end}}}},
			},
		},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
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
	// lambda.Start(Handler)
	Handler(context.TODO())
}
