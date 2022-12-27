package hanlders

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/types"
)

type EditNotificationHandler struct {
	Body types.TelegramMessageWrapper
}

func (handler EditNotificationHandler) Handle() error {
	requestBody := handler.Body
	tgRequestTemplate := "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

	env := config.GetInstance()
	botId := env.SwimNotifyBotId
	messageText := requestBody.Message.Text
	chatId := requestBody.Message.Chat.Id
	re := regexp.MustCompile(`^/editmsg(@SwimNotifyBot)? (.*)`)
	matched := re.FindStringSubmatch(messageText)

	if len(matched) == 0 {
		requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "empty msg is not allowed")
		helper.SendGetRequest(requestURL)

		return fmt.Errorf("failed")
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
		return fmt.Errorf("failed")
	}

	requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "done")
	helper.SendGetRequest(requestURL)
	requestURL = fmt.Sprintf(tgRequestTemplate, botId, chatId, msgForStore)
	helper.SendGetRequest(requestURL)

	return nil
}
