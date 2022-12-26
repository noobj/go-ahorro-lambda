package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
)

type TelegramMessageWrapper struct {
	Message  TelegramMessageBody `json:"message"`
	UpdateId string              `json:"update_id"`
}

type TelegramMessageBody struct {
	Chat TelegramMessageChat `json:"chat"`
	Text string              `json:"text"`
}

type TelegramMessageChat struct {
	Id int64 `json:"id"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	env := config.GetInstance()
	botId := env.SwimNotifyBotId
	tgRequestTemplate := "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"
	fmt.Printf("%+v", request)
	var body TelegramMessageWrapper
	json.Unmarshal([]byte(request.Body), &body)
	fmt.Printf("%+v", body)
	messageText := body.Message.Text
	chatId := body.Message.Chat.Id

	if matched, _ := regexp.MatchString("^/editmsg(@SwimNotifyBot)?", messageText); !matched {
		requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "shut up")
		helper.SendGetRequest(requestURL)

		return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("failed")
	}

	re := regexp.MustCompile(`^/editmsg(@SwimNotifyBot)? (.*)`)
	matched := re.FindStringSubmatch(messageText)

	if len(matched) == 0 {
		requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "empty msg is not allowed")
		helper.SendGetRequest(requestURL)

		return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("failed")
	}

	msgForStore := matched[len(matched)-1]

	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	item := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String("SWIM"),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":msg": {S: aws.String(msgForStore)},
		},
		UpdateExpression: aws.String("SET Msg=:msg"),
		TableName:        aws.String(env.DynamoSwimbotMsgTable),
	}

	_, err := svc.UpdateItem(item)
	if err != nil {
		log.Println("Task is done, but update status failed", err)

		requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "failed")
		helper.SendGetRequest(requestURL)
		return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("failed")
	}

	requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "done")
	helper.SendGetRequest(requestURL)

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("done")
}

func main() {
	lambda.Start(Handler)
}
