package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	env := config.GetInstance()
	botId := env.SwimNotifyBotId
	channelId := env.SwimNotifyChannelId
	message, messageExist := request.QueryStringParameters["message"]

	session, _ := session.NewSession()
	svc := dynamodb.New(session)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String("SWIM"),
			},
		},
		TableName: aws.String(env.DynamoSwimbotMsgTable),
	}

	item, err := svc.GetItem(input)

	defaultContent := "[溫腥提醒]各位奴才們，明天又到了一週最開心的週二看妹日囉😍，請別忘了 攜帶泳具👙，喵~"

	msgToSend := ""

	if err != nil {
		msgToSend = defaultContent
		fmt.Println(err.Error())
	} else {
		msgToSend = *item.Item["Msg"].S
	}

	content := url.QueryEscape(msgToSend)

	if messageExist {
		content = url.QueryEscape(message)
	}

	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botId, channelId, content)

	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return helper.GenerateErrorResponse[events.APIGatewayProxyResponse](401)
	}
	log.Println(res)
	if res.StatusCode != 200 {
		body, error := ioutil.ReadAll(res.Body)
		res.Body.Close()
		log.Panicln(string(body), error)
	}

	return helper.GenerateApiResponse[events.APIGatewayProxyResponse]("sent")
}

func main() {
	lambda.Start(Handler)
}
