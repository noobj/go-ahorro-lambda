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

type Cate struct {
	Category primitive.ObjectID `bson:"category,omitempty"`
}

type EntryCatgegoryBundle struct {
	Entries  []EntryRepository.Entry    `bson:"entries,omitempty"`
	Sum      int                        `bson:"sum,omitempty"`
	Category []EntryRepository.Category `bson:"category,omitempty"`
}

func checkTimeFormat(format string, timeString string) bool {
	_, err := time.Parse(format, timeString)

	return err == nil
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

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "$expr", Value: bson.D{
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$user", userId}}},
				bson.D{{Key: "$gte", Value: bson.A{"$date", startFromQuery}}},
				bson.D{{Key: "$lte", Value: bson.A{"$date", endFromQuery}}},
				//TODO: how do I spread the slice?
				// excludeCondition
			},
			}},
		}},
	}}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: sortColumn, Value: -1}},
	}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$category"},
		{Key: "entries", Value: bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "amount", Value: "$amount"},
				{Key: "date", Value: "$date"},
				{Key: "descr", Value: "$descr"},
			}}},
		},
		{Key: "sum", Value: bson.D{{
			Key: "$sum", Value: "$amount"},
		}},
	},
	}}

	sortSumStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: "sum", Value: -1}},
	}}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "categories"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "category"},
	},
	}}

	repoResults := entryRepository.Aggregate([]bson.D{matchStage, sortStage, groupStage, sortSumStage, lookupStage})

	for repoResult := range repoResults {
		doc, _ := bson.Marshal(repoResult)

		var test EntryCatgegoryBundle
		bson.Unmarshal(doc, &test)

		fmt.Printf("%+v\n", test)
	}

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
