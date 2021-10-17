package runtime

type TokenProxy interface {
	GetChatId() uint
	GetState() string
	GetScenarioName() string
	GetIsBlocked() bool
	GetTimeOffset() int
	GetUserName() string
	GetFirstName() string
	GetLastName() string
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
