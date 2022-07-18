package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	container "github.com/golobby/container/v3"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	EntryRepository "github.com/noobj/swim-crowd-lambda-go/internal/repositories/ahorro"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response events.APIGatewayProxyResponse

func checkTimeFormat(format string, timeString string) bool {
	_, err := time.Parse(format, timeString)

	if err != nil {
		return false
	}

	return true
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	startFromQuery, startExist := request.QueryStringParameters["timeStart"]
	endFromQuery, endExist := request.QueryStringParameters["timeEnd"]
	// categoriesExcludeInput, cateExcludeExist := request.QueryStringParameters["categoriesExclude"]
	sortByDateInput, sortExist := request.QueryStringParameters["entriesSortByDate"]

	if !checkTimeFormat("2006-01-02", startFromQuery) || !checkTimeFormat("2006-01-02", endFromQuery) || !startExist || !endExist {
		panic("something wrong with time query string")
	}

	var sortColumn string

	if sortExist && sortByDateInput == "true" {
		sortColumn = "date"
	} else {
		sortColumn = "amount"
	}

	// excludeCondition := []bson.D{}

	// if cateExcludeExist {
	// 	for _, cate := range strings.Split(categoriesExcludeInput, ",") {
	// 		cateId, _ := primitive.ObjectIDFromHex(cate)
	// 		condition := bson.D{{"$ne", bson.A{"$category", cateId}}}
	// 		excludeCondition = append(excludeCondition, condition)
	// 	}
	// }
	var entryRepository repositories.Repository
	container.Resolve(&entryRepository)
	defer entryRepository.Disconnect()()

	userId, _ := primitive.ObjectIDFromHex("627106d67b2f25ddd3daf964")

	matchStage := bson.D{{"$match", bson.D{
		{"$expr", bson.D{
			{"$and", bson.A{
				bson.D{{"$eq", bson.A{"$user", userId}}},
				bson.D{{"$gte", bson.A{"$date", "2022-01-01"}}},
				bson.D{{"$lte", bson.A{"$date", "2022-01-11"}}},
				//TODO: how do I spread the slice?
				// excludeCondition
			},
			}},
		}},
	}}

	sortStage := bson.D{{"$sort", bson.D{
		{sortColumn, -1}},
	}}

	repoResults := entryRepository.Aggregate([]bson.D{matchStage, sortStage})

	fmt.Println(repoResults)

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}

	return resp, nil
}

func main() {
	container.Singleton(func() repositories.Repository {
		return EntryRepository.New()
	})

	lambda.Start(Handler)
}
