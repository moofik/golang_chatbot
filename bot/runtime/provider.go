package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/config"
	"bot-daedalus/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ChatProvider interface {
	GetCommand() command.Command
	GetToken() *models.Token
	GetConfig() config.ProviderConfig
	SendTextMessage(text string) error
}

type TelegramProvider struct {
	config  config.ProviderConfig
	message Message
}

func (p *TelegramProvider) GetCommand() command.Command {
	m := p.GetTypedMessage(p.message)

	fmt.Printf("%v\n\n", m)

	if m.Message.Text[0] == '/' {
		return &command.UserInputCommand{Text: m.Message.Text}
	}

	return &command.MockCommand{}
}

func (p *TelegramProvider) getTokedId() uint {
	return p.GetTypedMessage(p.message).Message.Chat.Id
}

func (p *TelegramProvider) GetToken() *models.Token {
	//find token in DB or create new
	msg := p.GetTypedMessage(p.message)
	return &models.Token{
		ChatId: p.getTokedId(),
		FirstName: msg.Message.Chat.FirstName,
		LastName: msg.Message.Chat.LastName,
	}
}

func (p *TelegramProvider) GetConfig() config.ProviderConfig {
	return p.config
}

func (p *TelegramProvider) GetTypedMessage(m Message) TelegramMessage {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(m)

	var telegramMessage TelegramMessage
	_ = json.Unmarshal([]byte(buf.String()), &telegramMessage)

	return telegramMessage
}

func (p *TelegramProvider) SendTextMessage(text string) error {
	reqBody := &OutgoingMessage{
		ChatID: p.GetToken().ChatId,
		Text:   text,
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