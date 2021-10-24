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
	DB *gorm.DB
}

func (tf TokenFactory) GetOrCreate(p runtime.ChatProvider) runtime.TokenProxy {
	repository := &TokenRepository{DB: tf.DB}
	var token runtime.TokenProxy

	if p.GetConfig().Name == runtime.PROVIDER_TELEGRAM {
		msg := runtime.GetTelegramMessage(p.GetMessage())

		if msg.Message.Chat.Id != 0 {
			token = repository.FindByChatIdAndScenario(int(msg.Message.Chat.Id), p.GetScenarioName())

			if token == nil {
				token = &Token{
					ChatId:       msg.Message.Chat.Id,
					UserName:     msg.Message.Chat.UserName,
					FirstName:    msg.Message.Chat.FirstName,
					LastName:     msg.Message.Chat.LastName,
					ScenarioName: p.GetScenarioName(),
					State:        "unknown",
				}
			}
		} else if msg.CallbackQuery.Chat.Id != 0 {
			token := repository.FindByChatIdAndScenario(int(msg.CallbackQuery.Chat.Id), p.GetScenarioName())

			if token == nil {
				token = &Token{
					ChatId:       msg.CallbackQuery.Chat.Id,
					UserName:     msg.CallbackQuery.Chat.UserName,
					FirstName:    msg.CallbackQuery.Chat.FirstName,
					LastName:     msg.CallbackQuery.Chat.LastName,
					ScenarioName: p.GetScenarioName(),
					State:        "unknown",
				}
			}
		}

		return token
	}

	panic("Token not found")
}

type TokenRepository struct {
	DB *gorm.DB
}

func (r TokenRepository) FindByChatIdAndScenario(chatId int, scenario string) runtime.TokenProxy {
	var token Token
	res := r.DB.First(&token, "chat_id = ? and scenario_name = ?", chatId, scenario)
	if res.Error != nil {
		return nil
	}
	return &token
}

func (r TokenRepository) Persist(token runtime.TokenProxy) {
	r.DB.Save(token)
}

func (r TokenRepository) FindByScenario(scenario string) []runtime.TokenProxy {
	var tokens []Token
	r.DB.Where("scenario_name = ?", scenario).Find(&tokens)

	models := make([]runtime.TokenProxy, len(tokens))
	for i, v := range tokens {
		models[i] = runtime.TokenProxy(v)
	}

	return models
}

func isNil(i runtime.TokenProxy) bool {
	var ret bool
	switch i.(type) {
	case *Token:
		v := i.(*Token)
		ret = v == nil
	default:
		ret = false
	}
	return ret
}
