package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly"
	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	fmt.Printf(fmt.Sprintf("Scrapping at %s", time.Now().Format("2006-01-02 15:04")))

	var crowdCounts []int

	c := colly.NewCollector(
		colly.AllowedDomains("tndcsc.com.tw"),
	)

	c.OnHTML(".w3_agile_logo", func(e *colly.HTMLElement) {
		count, _ := strconv.Atoi(e.ChildText("p font"))
		crowdCounts = append(crowdCounts, count)
	})

	c.Visit("http://tndcsc.com.tw/index.aspx")

	c.Wait()

	client := mongodb.ClientInstance
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("swimCrowdDB").Collection("entries")
	doc := bson.D{{"amount", crowdCounts[2]}, {"time", time.Now().Format("2006-01-02 15:04")}}
	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Scrapping done.")

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
