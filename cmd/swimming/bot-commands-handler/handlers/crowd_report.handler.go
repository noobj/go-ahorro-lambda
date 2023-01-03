package hanlders

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/types"
)

type CrowdReportHandler struct {
	Body types.TelegramMessageWrapper
}

func (handler CrowdReportHandler) Handle() error {
	requestBody := handler.Body
	tgRequestTemplate := "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

	env := config.GetInstance()
	botId := env.SwimNotifyBotId
	chatId := requestBody.Message.Chat.Id
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

	msgToSend := fmt.Sprintf("目前有 %d 位辣妹在泳池👙", crowdCounts[2])
	requestURL := fmt.Sprintf(tgRequestTemplate, botId, chatId, msgToSend)
	helper.SendGetRequest(requestURL)

	return nil
}
