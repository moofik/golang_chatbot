package runtime

type TokenProxy interface {
	GetId() uint
	GetChatId() uint
	GetState() string
	GetScenarioName() string
	GetIsBlocked() bool
	GetTimeOffset() int
	GetUserName() string
	GetFirstName() string
	GetLastName() string
	GetExtras() map[string]string
	SetExtras(extras map[string]string)
	GetLastBotMessageId() int
	SetLastBotMessageId(id int)
	GetIsLastBotMessageRemovable() bool
	SetIsLastBotMessageRemovable(status bool)
	ToPlainStruct() interface{}
}

type TokenFactory interface {
	GetOrCreate(p ChatProvider) TokenProxy
}

type TokenRepository interface {
	FindByChatIdAndScenario(chatId int, scenario string) TokenProxy
	Persist(token TokenProxy)
	FindByScenario(scenario string) []TokenProxy
}
