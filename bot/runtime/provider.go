package runtime

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/config"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const PROVIDER_TELEGRAM string = "telegram"
const PROVIDER_VK string = "vk"
const PROVIDER_WHATSAPP string = "whatsapp"
const PROVIDER_FACEBOOK string = "facebook"

type ChatProvider interface {
	GetCommand(state *State) command.Command
	GetMessageFactory() SerializedMessageFactory
	GetMessage() Message
	GetToken() TokenProxy
	GetConfig() config.ProviderConfig
	GetScenarioName() string
	GetTokenRepository() TokenRepository
	SendTextMessage(text string, ctx ProviderContext) error
}

type TelegramOutgoingMessage struct {
	ChatID      uint                             `json:"chat_id"`
	Text        string                           `json:"text"`
	ParseMode   string                           `json:"parse_mode"`
	ReplyMarkup map[string][][]map[string]string `json:"reply_markup,omitempty"`
}

type TelegramProvider struct {
	tokenFactory   TokenFactory
	scenarioName   string
	messageFactory SerializedMessageFactory
	config         config.ProviderConfig
	message        Message
	TokenRepository
}

func (p *TelegramProvider) GetTokenRepository() TokenRepository {
	return p.TokenRepository
}

func (p *TelegramProvider) GetCommand(state *State) command.Command {
	m := GetTelegramMessage(p.message)

	if m.Message.Text == "" && m.CallbackQuery.Data == "" {
		return nil
	}

	cmdFlag := false
	cmdText := ""

	if m.Message.Text != "" && m.Message.Text[0] == '/' {
		cmdFlag = true
		cmdText = m.Message.Text
	} else if m.CallbackQuery.Data != "" && m.CallbackQuery.Data[0] == '/' {
		cmdFlag = true
		cmdText = m.CallbackQuery.Data
	}

	if cmdFlag {
		var dataSlice []string = []string{cmdText}
		var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
		for i, d := range dataSlice {
			interfaceSlice[i] = d
		}
		return command.CreateCommand("button", state.Name, interfaceSlice)
	}

	var dataSlice []string = []string{m.Message.Text}
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return command.CreateCommand("text_input", state.Name, interfaceSlice)
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

func (p *TelegramProvider) GetScenarioName() string {
	return p.scenarioName
}

func (p *TelegramProvider) GetMessage() Message {
	return p.message
}

func (p *TelegramProvider) SendTextMessage(text string, ctx ProviderContext) error {
	buttons := ctx.State.TransitionStorage.AllButtonCommands()

	var buttonsSlice []map[string]string

	for _, button := range buttons {
		buttonsSlice = append(buttonsSlice, map[string]string{
			"text":          button.GetCaption(),
			"callback_data": button.GetInput(),
		})
	}

	reqBody := &TelegramOutgoingMessage{
		ChatID:    p.GetToken().GetChatId(),
		Text:      text,
		ParseMode: "HTML",
	}

	if len(buttonsSlice) > 0 {
		reqBody.ReplyMarkup = map[string][][]map[string]string{
			"inline_keyboard": {buttonsSlice},
		}
	}

	reqBytes, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	url := "https://api.telegram.org/bot" + p.GetConfig().Token + "/sendMessage"

	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

type ProviderContext struct {
	State   *State
	Command command.Command
	Token   TokenProxy
}
