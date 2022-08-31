package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golobby/container/v3"
	"github.com/mitchellh/mapstructure"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/mongodb"
	"github.com/noobj/go-serverless-services/internal/repositories"
	AhorroRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro"
	CategoryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/category"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type AhorroRes struct {
	Tables []AhorroJsonFormat
}

type AhorroJsonFormat struct {
	TableName    string
	Items        []map[string]interface{}
	RowCounts    int
	ColumnsCount int
	ColumnNames  interface{}
}

type EntryItem struct {
	Amount     string
	Date       string
	Descr      string
	CategoryId string `mapstructure:"category_id"`
}

type CategoryItem struct {
	Id   string `mapstructure:"_id"`
	Name string
}

func Handler(ctx context.Context, event events.SQSEvent) (events.APIGatewayProxyResponse, error) {
	var userRepository repositories.IRepository
	container.NamedResolve(&userRepository, "UserRepo")
	var categoryRepository repositories.IRepository
	container.NamedResolve(&categoryRepository, "CategoryRepo")
	var entryRepository repositories.IRepository
	container.NamedResolve(&entryRepository, "EntryRepo")

	userId := *event.Records[0].MessageAttributes["UserId"].StringValue
	var user UserRepository.User
	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	err := userRepository.FindOne(context.TODO(), bson.M{"_id": userObjectId}).Decode(&user)
	if err != nil {
		log.Println(err)
		return helper.GenerateInternalErrorResponse("error: user not found")
	}
	config := helper.GenerateOauthConfig()

	token := oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  user.GoogleAccessToken,
		RefreshToken: user.GoogleRefreshToken,
		// TODO: use real value
		Expiry: time.Now().Add(time.Hour * -2),
	}

	client := config.Client(ctx, &token)
	service, _ := drive.New(client)

	file, err := service.Files.List().Q("name contains 'ahorro'").OrderBy("createdTime desc").PageSize(1).Do()

	if err != nil {
		log.Printf("google service error: %v", err)
		return helper.GenerateInternalErrorResponse("error: google service error")
	}

	fileId := file.Files[0].Id

	content, err := service.Files.Get(fileId).Download()

	if err != nil {
		log.Printf("google service error: %v", err)
		return helper.GenerateInternalErrorResponse("error: google service error")
	}

	buff := make([]byte, 10)
	var tmp []string

	for {
		n, err := content.Body.Read(buff)
		if err == io.EOF {
			break
		}

		tmp = append(tmp, string(buff[:n]))
	}

	var ahorroRes AhorroRes
	json.Unmarshal([]byte(strings.Join(tmp, "")), &ahorroRes)
	var entryItems []EntryItem
	var categoryItems []CategoryItem
	for _, table := range ahorroRes.Tables {
		if table.TableName == "expense" {
			err = mapstructure.Decode(table.Items, &entryItems)
			if err != nil {
				log.Println("Json format error in the expense sector", err)
				return helper.GenerateInternalErrorResponse()
			}
		}

		if table.TableName == "category" {
			err = mapstructure.Decode(table.Items, &categoryItems)
			if err != nil {
				log.Println("Json format error in the category sector", err)
				return helper.GenerateInternalErrorResponse()
			}
		}
	}

	err = mongodb.GetInstance().UseSession(ctx, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction()
		if err != nil {
			return err
		}

		_, err = categoryRepository.DeleteMany(sc, bson.M{"user": userObjectId})
		if err != nil {
			return err
		}

		categoriesForInsert, newCategoryIdMap := collateCategoryItems(categoryItems, userObjectId)

		if _, err = categoryRepository.InsertMany(sc, categoriesForInsert); err != nil {
			return err
		}

		if _, err = entryRepository.DeleteMany(sc, bson.M{"user": userObjectId}); err != nil {
			return err
		}

		entriesForInsert := collateEntryItems(entryItems, newCategoryIdMap, userObjectId)
		_, err = entryRepository.InsertMany(sc, entriesForInsert)

		if err != nil {
			return err
		}

		if err = sc.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("Something wrong in Transaction", err)
		return helper.GenerateInternalErrorResponse()
	}

	return helper.GenerateApiResponse("ok")
}

func collateCategoryItems(categoryItems []CategoryItem, userId primitive.ObjectID) ([]interface{}, map[string]primitive.ObjectID) {
	newCategoryIdMap := make(map[string]primitive.ObjectID)
	var result []interface{}
	for _, categoryItem := range categoryItems {
		if BuiltinCategories[categoryItem.Name] != "" {
			categoryItem.Name = BuiltinCategories[categoryItem.Name]
		}

		rand.Seed(time.Now().UnixNano())
		colorString := fmt.Sprintf("#%06x", rand.Intn(16777215))
		newId := primitive.NewObjectID()
		newItemBson := bson.M{
			"_id":   newId,
			"name":  categoryItem.Name,
			"user":  userId,
			"color": colorString,
		}
		result = append(result, newItemBson)

		newCategoryIdMap[categoryItem.Id] = newId
	}

	return result, newCategoryIdMap
}

func collateEntryItems(entryItems []EntryItem, cateIdMap map[string]primitive.ObjectID, userId primitive.ObjectID) []interface{} {
	var result []interface{}
	for _, entryItem := range entryItems {
		newItemBson := bson.M{
			"amount":   entryItem.Amount,
			"date":     entryItem.Date,
			"descr":    entryItem.Descr,
			"category": cateIdMap[entryItem.CategoryId],
			"user":     userId,
		}
		result = append(result, newItemBson)
	}

	return result
}

func main() {
	userRepo := UserRepository.New()
	defer userRepo.Disconnect()()
	container.NamedSingleton("UserRepo", func() repositories.IRepository {
		return userRepo
	})

	container.NamedSingleton("CategoryRepo", func() repositories.IRepository {
		return CategoryRepository.New()
	})

	container.NamedSingleton("EntryRepo", func() repositories.IRepository {
		return AhorroRepository.New()
	})

	lambda.Start(Handler)
}
