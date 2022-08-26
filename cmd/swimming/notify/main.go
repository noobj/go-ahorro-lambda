package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	botId := os.Getenv("SWIM_NOTIFY_BOT_ID")
	channelId := os.Getenv("SWIM_NOTIFY_CHANNEL_ID")

	content := url.QueryEscape("Today is the day for Swimminggg🌊💪🏽\nDon't forget to bring the gears")
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botId, channelId, content)

	_, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return helper.GenerateInternalErrorResponse()
	}

	return helper.GenerateApiResponse("sent")
}

func main() {
	lambda.Start(Handler)
}
