package runtime

type Setting interface {
	IsOffline() bool
	SetOffline(offline bool)
	IsMaintenance() bool
	SetMaintenance(maintenance bool)
	GetTelegramAdminsIds() []int
	SetTelegramAdminsIds(ids []int)
	GetTelegramNotificationChannelsTokens() []string
	SetTelegramNotificationChannelsTokens(ids []string)
}

type SettingsRepository interface {
	FindByScenarioName(scenarioName string) Setting
	Persist(setting Setting)
}
