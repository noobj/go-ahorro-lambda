package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/config"
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
				env := config.GetInstance()
				mongoUser := env.MongoUser
				mongoPassword := env.MongoPassword
				mongoPath := env.MongoPath
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
