package models

import (
	"bot-daedalus/bot/runtime"
	"encoding/json"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	ChatId                    uint
	State                     string
	ScenarioName              string
	IsBlocked                 bool
	TimeOffset                int
	UserName                  string
	FirstName                 string
	LastName                  string
	Extras                    string
	LastBotMessageId          int
	IsLastBotMessageRemovable bool
}

func (t *Token) GetExtras() map[string]string {
	var objmap map[string]string
	err := json.Unmarshal([]byte(t.Extras), &objmap)

	if err != nil {
		panic(err)
	}

	return objmap
}

func (t *Token) SetExtras(extras map[string]string) {
	result, err := json.Marshal(extras)

	if err != nil {
		panic(err)
	}

	t.Extras = string(result)
}

func (t *Token) GetId() uint {
	return t.ID
}

func (t *Token) GetChatId() uint {
	return t.ChatId
}

func (t *Token) GetState() string {
	return t.State
}

func (t *Token) GetScenarioName() string {
	return t.ScenarioName
}

func (t *Token) GetIsBlocked() bool {
	return t.IsBlocked
}

func (t *Token) GetTimeOffset() int {
	return t.TimeOffset
}

func (t *Token) GetUserName() string {
	return t.UserName
}

func (t *Token) GetFirstName() string {
	return t.FirstName
}

func (t *Token) GetLastName() string {
	return t.LastName
}

func (t *Token) GetLastBotMessageId() int {
	return t.LastBotMessageId
}

func (t *Token) SetLastBotMessageId(id int) {
	t.LastBotMessageId = id
}

func (t *Token) ToPlainStruct() interface{} {
	return t
}

func (t *Token) GetIsLastBotMessageRemovable() bool {
	return t.IsLastBotMessageRemovable
}

func (t *Token) SetIsLastBotMessageRemovable(status bool) {
	t.IsLastBotMessageRemovable = status
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
					Extras:       "{}",
				}
			}
		} else if msg.CallbackQuery.QueryId != "" {
			token = repository.FindByChatIdAndScenario(int(msg.CallbackQuery.From.Id), p.GetScenarioName())

			if token == nil {
				token = &Token{
					ChatId:       msg.CallbackQuery.From.Id,
					UserName:     msg.CallbackQuery.From.UserName,
					FirstName:    msg.CallbackQuery.From.FirstName,
					LastName:     msg.CallbackQuery.From.LastName,
					ScenarioName: p.GetScenarioName(),
					State:        "unknown",
					Extras:       "{}",
				}
			}
		} else {
			panic(msg)
		}

		return token
	}

	panic("Token not found")
}

type TokenRepository struct {
	DB *gorm.DB
}

func (r *TokenRepository) FindById(id int) runtime.TokenProxy {
	var token Token
	res := r.DB.First(&token, "id = ?", id)
	if res.Error != nil {
		return nil
	}
	return &token
}

func (r *TokenRepository) FindByChatIdAndScenario(chatId int, scenario string) runtime.TokenProxy {
	var token Token
	res := r.DB.First(&token, "chat_id = ? and scenario_name = ?", chatId, scenario)
	if res.Error != nil {
		return nil
	}
	return &token
}

func (r *TokenRepository) Persist(token runtime.TokenProxy) {
	r.DB.Save(token)
}

func (r *TokenRepository) FindByScenario(scenario string) []runtime.TokenProxy {
	var tokens []*Token
	r.DB.Where("scenario_name = ?", scenario).Find(&tokens)

	models := make([]runtime.TokenProxy, len(tokens))
	for i, v := range tokens {
		models[i] = runtime.TokenProxy(v)
	}

	return models
}

func (r *TokenRepository) FindByScenarioIdsNotIn(scenario string, listIds []uint) []runtime.TokenProxy {
	var tokens []*Token

	r.DB.
		Where("scenario_name = ?", scenario).
		Not(map[string]interface{}{"chat_id": listIds}).
		Find(&tokens)

	models := make([]runtime.TokenProxy, len(tokens))
	for i, v := range tokens {
		models[i] = runtime.TokenProxy(v)
	}

	return models
}
