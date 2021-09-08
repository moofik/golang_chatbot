package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const PROVIDER_TELEGRAM string = "teleram"
const PROVIDER_VK string = "vk"
const PROVIDER_WHATSAPP string = "whatsapp"
const PROVIDER_FACEBOOK string = "facebook"

type ChatProvider interface {
	GetCommand() command.Command
	GetMessageFactory() SerializedMessageFactory
	GetToken() TokenProxy
	GetConfig() config.ProviderConfig
	SendTextMessage(text string) error
}

type OutgoingMessage struct {
	ChatID    uint   `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type TelegramProvider struct {
	tokenFactory   TokenFactory
	messageFactory SerializedMessageFactory
	config         config.ProviderConfig
	message        Message
}

func (p *TelegramProvider) GetCommand() command.Command {
	m := GetTelegramMessage(p.message)

	fmt.Printf("%v\n\n", m)

	if m.Message.Text[0] == '/' {
		return &command.UserInputCommand{Text: m.Message.Text}
	}

	return &command.MockCommand{}
}

func (p *TelegramProvider) getTokedId() uint {
	return GetTelegramMessage(p.message).Message.Chat.Id
}

func (p *TelegramProvider) GetMessageFactory() SerializedMessageFactory {
	return p.messageFactory
}

func (p *TelegramProvider) GetToken() TokenProxy {
	//find token in DB or create new
	return p.tokenFactory.GetOrCreate(p)
}

func (p *TelegramProvider) GetConfig() config.ProviderConfig {
	return p.config
}

func (p *TelegramProvider) SendTextMessage(text string) error {
	reqBody := &OutgoingMessage{
		ChatID:    p.GetToken().GetChatId(),
		Text:      text,
		ParseMode: "HTML",
	}

	reqBytes, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	fmt.Println("Run send text message action")
	url := "https://api.telegram.org/bot" + p.GetConfig().Token + "/sendMessage"
	fmt.Println(url)
	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		panic("ERR")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
