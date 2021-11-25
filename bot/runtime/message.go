package runtime

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Message interface{}

type TelegramMessage struct {
	UpdateId uint `json:"update_id,omitempty"`
	Message  struct {
		MessageId uint `json:"message_id,omitempty"`
		From      struct {
			Id           uint   `json:"id,omitempty"`
			IsBot        bool   `json:"is_bot,omitempty"`
			FirstName    string `json:"first_name,omitempty"`
			LastName     string `json:"last_name,omitempty"`
			UserName     string `json:"username,omitempty"`
			LanguageCode string `json:"language_code,omitempty"`
		} `json:"from,omitempty"`
		Chat struct {
			Id        uint   `json:"id,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			UserName  string `json:"username,omitempty"`
			Type      string `json:"type,omitempty"`
		} `json:"chat,omitempty"`
		Date uint   `json:"date,omitempty"`
		Text string `json:"text,omitempty"`
	} `json:"message,omitempty"`
	CallbackQuery struct {
		QueryId string `json:"id,omitempty"`
		From    struct {
			Id           uint   `json:"id,omitempty"`
			IsBot        bool   `json:"is_bot,omitempty"`
			FirstName    string `json:"first_name,omitempty"`
			LastName     string `json:"last_name,omitempty"`
			UserName     string `json:"username,omitempty"`
			LanguageCode string `json:"language_code,omitempty"`
		} `json:"from,omitempty"`
		Data string `json:"data,omitempty"`
	} `json:"callback_query,omitempty"`
}

type TelegramOutgoingMessage struct {
	ChatID      uint                             `json:"chat_id"`
	Text        string                           `json:"text"`
	ParseMode   string                           `json:"parse_mode"`
	ReplyMarkup map[string][][]map[string]string `json:"reply_markup,omitempty"`
}

type TelegramOutgoingDeleteMessage struct {
	ChatID    uint `json:"chat_id"`
	MessageID uint `json:"message_id"`
}

type SerializedMessageFactory interface {
	GetSerializedMessage(c ProviderConfig) Message
}

type DefaultSerializedMessageFactory struct {
	Ctx *gin.Context
}

func (f *DefaultSerializedMessageFactory) GetSerializedMessage(c ProviderConfig) Message {
	if c.Name == "telegram" {
		var json TelegramMessage

		if f.Ctx.BindJSON(&json) == nil {
			return json
		}

		return nil
	}

	panic("unknown message type")
}

func GetTelegramMessage(m Message) TelegramMessage {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(m)

	var telegramMessage TelegramMessage
	_ = json.Unmarshal([]byte(buf.String()), &telegramMessage)

	return telegramMessage
}
