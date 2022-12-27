package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	hanlders "github.com/noobj/go-serverless-services/cmd/swimming/bot-commands-handler/handlers"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/types"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	env := config.GetInstance()
	botId := env.SwimNotifyBotId
	tgRequestTemplate := "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"
	fmt.Printf("%+v", request)
	var requestBody types.TelegramMessageWrapper
	json.Unmarshal([]byte(request.Body), &requestBody)
	fmt.Printf("%+v", requestBody)

	messageText := requestBody.Message.Text
	chatId := requestBody.Message.Chat.Id

	if matched, _ := regexp.MatchString("^/editmsg(@SwimNotifyBot)?", messageText); !matched {
		requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, "shut up")
		helper.SendGetRequest(requestURL)

		return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("failed")
	}

	var handler hanlders.IHandler
	handler = hanlders.EditNotificationHandler{
		Body: requestBody,
	}

	if err := handler.Handle(); err != nil {
		return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("failed")
	}

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("done")
}

func main() {
	lambda.Start(Handler)
}
