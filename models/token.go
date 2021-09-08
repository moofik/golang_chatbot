package models

import (
	"bot-daedalus/bot/runtime"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	ChatId       uint
	State        string
	ScenarioName string
	IsBlocked    bool
	TimeOffset   int
	UserName     string
	FirstName    string
	LastName     string
}

func (t Token) GetChatId() uint {
	return t.ChatId
}

func (t Token) GetState() string {
	return t.State
}
func (t Token) GetScenarioName() string {
	return t.ScenarioName
}

func (t Token) GetIsBlocked() bool {
	return t.IsBlocked
}

func (t Token) GetTimeOffset() int {
	return t.TimeOffset
}

func (t Token) GetUserName() string {
	return t.UserName
}

func (t Token) GetFirstName() string {
	return t.FirstName
}

func (t Token) GetLastName() string {
	return t.LastName
}

func (t Token) ToPlainStruct() interface{} {
	return &Token{
		ChatId:       t.GetChatId(),
		State:        t.GetState(),
		ScenarioName: t.GetScenarioName(),
		IsBlocked:    t.GetIsBlocked(),
		TimeOffset:   t.GetTimeOffset(),
		UserName:     t.GetUserName(),
		FirstName:    t.GetFirstName(),
		LastName:     t.GetLastName(),
	}
}

// TokenFactory implementation
type TokenFactory struct {
}

func (tf TokenFactory) GetOrCreate(p runtime.ChatProvider) runtime.TokenProxy {
	if p.GetConfig().Name == runtime.PROVIDER_TELEGRAM {
		msg := runtime.GetTelegramMessage(p.GetMessageFactory().GetSerializedMessage(p.GetConfig()))
		return &Token{
			ChatId:    msg.Message.Chat.Id,
			UserName:  msg.Message.Chat.UserName,
			FirstName: msg.Message.Chat.FirstName,
			LastName:  msg.Message.Chat.LastName,
		}
	}

	panic("Token not found")
}
