package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoPath := os.Getenv("MONGO_PATH")
	uri := fmt.Sprintf("mongodb+srv://%s:%s%s", mongoUser, mongoPassword, mongoPath)

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(cmdMonitor))
	defer client.Disconnect(context.TODO())

	createLoginInfosCreatedAtTTLIndex(client)
	createUsersAccountUniIndex(client)
	if err != nil {
		fmt.Println("mongo connect error", err)
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}
}

func createLoginInfosCreatedAtTTLIndex(client *mongo.Client) {
	refreshTokenExpireTime, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_TIME"))
	col := client.Database("ahorro").Collection("loginInfos")
	mod := mongo.IndexModel{
		Keys: bson.M{
			"createdAt": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32(refreshTokenExpireTime)),
	}

	ind, err := col.Indexes().CreateOne(context.TODO(), mod)

	// Check if the CreateOne() method returned any errors
	if err != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		// API call returns string of the index name
		fmt.Println("CreateOne() index:", ind)
		fmt.Println("CreateOne() type:", reflect.TypeOf(ind))
	}
}

func createUsersAccountUniIndex(client *mongo.Client) {
	col := client.Database("ahorro").Collection("users")
	mod := mongo.IndexModel{
		Keys: bson.M{
			"account": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	ind, err := col.Indexes().CreateOne(context.TODO(), mod)

	// Check if the CreateOne() method returned any errors
	if err != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		// API call returns string of the index name
		fmt.Println("CreateOne() index:", ind)
		fmt.Println("CreateOne() type:", reflect.TypeOf(ind))
	}
}
