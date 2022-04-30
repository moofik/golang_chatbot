package runtime

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
)

type Action interface {
	Run(p ChatProvider, t TokenProxy, s *State, prev *State, c Command) ActionError
	GetName() string
}

type SendTextMessage struct {
	params map[string]interface{}
}

func (a *SendTextMessage) GetName() string {
	return "send_text"
}

func (a *SendTextMessage) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	tmpl, err := template.New("test").Parse(a.params["text"].(string))
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer
	data := t.GetExtras()
	if err := tmpl.Execute(&tpl, data); err != nil {
		return &GenericActionError{InnerError: err}
	}

	result := tpl.String()
	lastBotMessageId := uint(t.GetLastBotMessageId())
	err = p.SendTextMessage(result, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	if clear, ok := a.params["clear_previous"]; ok && clear.(bool) && t.GetIsLastBotMessageRemovable() {
		DeleteMessage(t.GetChatId(), lastBotMessageId, p.GetConfig().Token)

		if removable, ok := a.params["removable"]; ok && !removable.(bool) {
			t.SetIsLastBotMessageRemovable(false)
		} else {
			t.SetIsLastBotMessageRemovable(true)
		}
	} else {
		t.SetIsLastBotMessageRemovable(true)
	}

	return nil
}

type RememberInput struct {
	params map[string]interface{}
}

func (a *RememberInput) GetName() string {
	return "remember_input"
}

func (a *RememberInput) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	extras := t.GetExtras()
	extras[a.params["var"].(string)] = c.GetInput()
	t.SetExtras(extras)
	return nil
}

type RememberVar struct {
	params map[string]interface{}
}

func (a *RememberVar) GetName() string {
	return "remember_var"
}

func (a *RememberVar) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	extras := t.GetExtras()
	extras[a.params["var"].(string)] = a.params["value"].(string)
	t.SetExtras(extras)
	return nil
}

type RememberCaption struct {
	params map[string]interface{}
}

func (a *RememberCaption) GetName() string {
	return "remember_input"
}

func (a *RememberCaption) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	newCmd := &ButtonPressedCommand{
		ButtonCommand: c.GetInput(),
		ButtonText:    c.GetCaption(),
		Metadata: &CommandMetadata{
			Cmd:        "button",
			Place:      s.Name,
			Uniqueness: c.GetInput(),
		}}

	cmd, _ := prev.GetCommandByUniqueness(newCmd)

	extras := t.GetExtras()
	extras[a.params["var"].(string)] = cmd.GetCaption()
	t.SetExtras(extras)
	return nil
}

type SendReplyMarkup struct {
	params map[string]interface{}
}

func (a *SendReplyMarkup) GetName() string {
	return "send_reply_markup"
}

func (a *SendReplyMarkup) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	tmpl, err := template.New("test").Parse(a.params["text"].(string))
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer

	data := t.GetExtras()

	if err := tmpl.Execute(&tpl, data); err != nil {
		return &GenericActionError{InnerError: err}
	}
	rawButtons := a.params["buttons"].([]interface{})
	buttons := make([]string, len(rawButtons))

	for _, button := range rawButtons {
		buttons = append(buttons, button.(string))
	}

	lastBotMessageId := uint(t.GetLastBotMessageId())
	result := tpl.String()
	err = p.SendMarkupMessage(buttons, result, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})

	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	if clear, ok := a.params["clear_previous"]; ok && clear.(bool) && t.GetIsLastBotMessageRemovable() {
		DeleteMessage(t.GetChatId(), lastBotMessageId, p.GetConfig().Token)

		if removable, ok := a.params["removable"]; ok && !removable.(bool) {
			t.SetIsLastBotMessageRemovable(false)
		} else {
			t.SetIsLastBotMessageRemovable(true)
		}
	} else {
		t.SetIsLastBotMessageRemovable(true)
	}

	t.SetIsLastBotMessageRemovable(false)

	return nil
}

type SendPhoto struct {
	Params map[string]interface{}
}

func (a *SendPhoto) GetName() string {
	return "send_photo"
}

func (a *SendPhoto) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	tmpl, err := template.New("test").Parse(a.Params["text"].(string))
	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer

	data := t.GetExtras()

	if err := tmpl.Execute(&tpl, data); err != nil {
		return &GenericActionError{InnerError: err}
	}

	var buttons []string

	if a.Params["buttons"] != nil {
		rawButtons := a.Params["buttons"].([]interface{})
		buttons = make([]string, len(rawButtons))

		for _, button := range rawButtons {
			buttons = append(buttons, button.(string))
		}
	}

	lastBotMessageId := uint(t.GetLastBotMessageId())
	result := tpl.String()

	var replyMarkup *TelegramReplyMarkup

	if a.Params["remove_keyboard"].(bool) == true {
		replyMarkup = &TelegramReplyMarkup{RemoveKeyboard: true}
	}

	err = p.SendLocalPhoto(buttons, result, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	}, replyMarkup)

	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	if clear, ok := a.Params["clear_previous"]; ok && clear.(bool) && t.GetIsLastBotMessageRemovable() {
		DeleteMessage(t.GetChatId(), lastBotMessageId, p.GetConfig().Token)

		if removable, ok := a.Params["removable"]; ok && !removable.(bool) {
			t.SetIsLastBotMessageRemovable(false)
		} else {
			t.SetIsLastBotMessageRemovable(true)
		}
	} else {
		t.SetIsLastBotMessageRemovable(true)
	}

	t.SetIsLastBotMessageRemovable(false)

	return nil
}

type SendChatAction struct {
	Params map[string]interface{}
}

func (a *SendChatAction) GetName() string {
	return "send_chat_action"
}

func (a *SendChatAction) Run(
	p ChatProvider,
	t TokenProxy,
	s *State,
	prev *State,
	c Command,
) ActionError {
	action := a.Params["action"].(string)
	lastBotMessageId := uint(t.GetLastBotMessageId())
	var replyMarkup *TelegramReplyMarkup

	if a.Params["remove_keyboard"].(bool) == true {
		replyMarkup = &TelegramReplyMarkup{RemoveKeyboard: a.Params["remove_keyboard"].(bool)}
	}

	err := p.SendChatAction(action, ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	}, replyMarkup)

	if err != nil {
		return &GenericActionError{InnerError: err}
	}

	if clear, ok := a.Params["clear_previous"]; ok && clear.(bool) && t.GetIsLastBotMessageRemovable() {
		DeleteMessage(t.GetChatId(), lastBotMessageId, p.GetConfig().Token)

		if removable, ok := a.Params["removable"]; ok && !removable.(bool) {
			t.SetIsLastBotMessageRemovable(false)
		} else {
			t.SetIsLastBotMessageRemovable(true)
		}
	} else {
		t.SetIsLastBotMessageRemovable(true)
	}

	t.SetIsLastBotMessageRemovable(false)

	return nil
}

func CreateAction(name string, params map[string]interface{}, actionRegistry func(string, map[string]interface{}) Action) Action {
	if name == "send_text" {
		return &SendTextMessage{params: params}
	}

	if name == "remember_input" {
		return &RememberInput{params: params}
	}

	if name == "remember_var" {
		return &RememberVar{params: params}
	}

	if name == "remember_caption" {
		return &RememberCaption{params: params}
	}

	if name == "send_reply_markup" {
		return &SendReplyMarkup{params: params}
	}

	if name == "send_photo" {
		return &SendPhoto{Params: params}
	}

	if name == "send_chat_action" {
		return &SendChatAction{Params: params}
	}

	if actionRegistry != nil {
		return actionRegistry(name, params)
	}

	return nil
}

type ActionError interface {
	error
}

type GenericActionError struct {
	InnerError error
}

func (m *GenericActionError) Error() string {
	return m.InnerError.Error()
}

func DeleteMessage(chatId, messageId uint, botToken string) {
	url := "https://api.telegram.org/bot" + botToken + "/deleteMessage"

	reqBody := &TelegramOutgoingDeleteMessage{
		ChatID:    uint(chatId),
		MessageID: uint(messageId),
	}

	reqBytes, err := json.Marshal(reqBody)

	_, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		panic(err)
	}
}
