package runtime

import (
	"bot-daedalus/config"
	"github.com/gin-gonic/gin"
)

type Message interface{}

type TelegramMessage struct {
	UpdateId uint `json:"update_id"`
	Message  struct {
		MessageId uint `json:"message_id"`
		From      struct {
			Id           uint   `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			UserName     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			Id        uint   `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date uint   `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

func GetMessage(ctx *gin.Context, c config.ProviderConfig) Message {
	if c.Name == "telegram" {
		var json TelegramMessage

		if ctx.BindJSON(&json) == nil {
			return json
		}

		return nil
	}

	panic("unknown message type")
}
