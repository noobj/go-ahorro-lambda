package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly"
	container "github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/repositories"
	EntryRepository "github.com/noobj/go-serverless-services/internal/repositories/swim/entry"
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

	var entryRepository repositories.IRepository
	container.Resolve(&entryRepository)

	doc := bson.D{{Key: "amount", Value: crowdCounts[2]}, {Key: "time", Value: time.Now().Format("2006-01-02 15:04")}}
	entryRepository.InsertOne(doc)

	fmt.Println("Scrapping done.")

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}

	return resp, nil
}

func main() {
	entryRepo := EntryRepository.New()
	container.Singleton(func() repositories.IRepository {
		return entryRepo
	})

	defer entryRepo.Disconnect()()

	lambda.Start(Handler)
}
