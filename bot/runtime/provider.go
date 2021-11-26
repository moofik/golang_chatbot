package runtime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const PROVIDER_TELEGRAM string = "telegram"
const PROVIDER_VK string = "vk"
const PROVIDER_WHATSAPP string = "whatsapp"
const PROVIDER_FACEBOOK string = "facebook"

type ChatProvider interface {
	GetCommand(state *State) Command
	GetMessageFactory() SerializedMessageFactory
	GetMessage() Message
	GetToken() TokenProxy
	GetConfig() ProviderConfig
	GetScenarioName() string
	GetTokenRepository() TokenRepository
	SendTextMessage(text string, ctx ProviderContext) error
	SendMarkupMessage(buttons []string, text string, ctx ProviderContext) error
}

type TelegramSendMessageResponse struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		MessageID int `json:"message_id,omitempty"`
		From      struct {
			ID        int    `json:"id,omitempty"`
			IsBot     bool   `json:"is_bot,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			Username  string `json:"username,omitempty"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id,omitempty"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Username  string `json:"username,omitempty"`
			Type      string `json:"type,omitempty"`
		} `json:"chat"`
		Date     int    `json:"date,omitempty"`
		Text     string `json:"text,omitempty"`
		Entities []struct {
			Offset int    `json:"offset,omitempty"`
			Length int    `json:"length,omitempty"`
			Type   string `json:"type,omitempty"`
		} `json:"entities"`
		ReplyMarkup struct {
			InlineKeyboard [][]struct {
				Text         string `json:"text,omitempty"`
				CallbackData string `json:"callback_data,omitempty"`
			} `json:"inline_keyboard,omitempty"`
		} `json:"reply_markup,omitempty"`
	} `json:"result,omitempty"`
}

type TelegramProvider struct {
	tokenFactory   TokenFactory
	scenarioName   string
	messageFactory SerializedMessageFactory
	config         ProviderConfig
	message        Message
	TokenRepository
}

func (p *TelegramProvider) GetTokenRepository() TokenRepository {
	return p.TokenRepository
}

func (p *TelegramProvider) GetCommand(state *State) Command {
	m := GetTelegramMessage(p.message)

	if m.Message.Text == "" && m.CallbackQuery.Data == "" {
		return nil
	}

	cmdFlag := false
	cmdText := ""

	if m.CallbackQuery.Data != "" && m.CallbackQuery.Data[0] == '/' {
		cmdFlag = true
		cmdText = m.CallbackQuery.Data
	}

	if cmdFlag {
		var dataSlice []string = []string{cmdText}
		var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
		for i, d := range dataSlice {
			interfaceSlice[i] = d
		}
		return CreateCommand("button", state.Name, interfaceSlice, nil)
	}

	var dataSlice []string = []string{m.Message.Text}
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
		interfaceSlice[i] = d
	}
	return CreateCommand("text_input", state.Name, interfaceSlice, nil)
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

func (p *TelegramProvider) GetConfig() ProviderConfig {
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

	var buttonsSlice [][]map[string]string

	for _, button := range buttons {
		buttonsSlice = append(buttonsSlice, []map[string]string{
			{
				"text":          button.GetCaption(),
				"callback_data": button.GetInput(),
			},
		})
	}

	reqBody := &TelegramOutgoingMessage{
		ChatID:    p.GetToken().GetChatId(),
		Text:      text,
		ParseMode: "HTML",
	}

	if len(buttonsSlice) > 0 {
		reqBody.ReplyMarkup = TelegramReplyMarkup{
			InlineKeyboard: buttonsSlice,
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
	defer res.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		sendMsgRes := GetTelegramSendMessageResponse(bodyString)
		ctx.Token.SetLastBotMessageId(sendMsgRes.Result.MessageID)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status" + res.Status)
	}

	return nil
}

func (p *TelegramProvider) SendMarkupMessage(buttons []string, text string, ctx ProviderContext) error {
	var buttonsSlice [][]map[string]string
	buttonsSlice = append(buttonsSlice, []map[string]string{})

	for _, button := range buttons {
		buttonsSlice[0] = append(buttonsSlice[0], map[string]string{
			"text": button,
		})
	}

	reqBody := &TelegramOutgoingMessage{
		ChatID:    p.GetToken().GetChatId(),
		Text:      text,
		ParseMode: "HTML",
	}

	if len(buttonsSlice) > 0 {
		reqBody.ReplyMarkup = TelegramReplyMarkup{
			Keyboard:       buttonsSlice,
			ResizeKeyboard: true,
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
	defer res.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		sendMsgRes := GetTelegramSendMessageResponse(bodyString)
		ctx.Token.SetLastBotMessageId(sendMsgRes.Result.MessageID)
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

type ProviderContext struct {
	State   *State
	Command Command
	Token   TokenProxy
}

func GetTelegramSendMessageResponse(body string) TelegramSendMessageResponse {
	var res TelegramSendMessageResponse
	_ = json.Unmarshal([]byte(body), &res)

	return res
}
