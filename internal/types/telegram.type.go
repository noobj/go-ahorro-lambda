package types

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
