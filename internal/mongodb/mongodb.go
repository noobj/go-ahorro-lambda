package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client

var once sync.Once

func GetInstance() *mongo.Client {
	if clientInstance == nil {
		once.Do(
			func() {
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

				if err != nil {
					panic(err)
				}
				fmt.Println("Creating single instance now.")
				clientInstance = client
			})
	}

	return clientInstance
}
