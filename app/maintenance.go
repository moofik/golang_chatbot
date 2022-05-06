package app

import (
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
)

type MaintenanceHandler struct {
	SettingsRepository *models.SettingsRepository
	ScenarioName       string
}

func (s *MaintenanceHandler) Handle(cmd runtime.Command, currentState *runtime.State, token runtime.TokenProxy) bool {
	setting := s.SettingsRepository.FindByScenarioName(s.ScenarioName)
	ids := setting.GetTelegramAdminsIds()
	chatId := token.GetChatId()
	isAdmin := false

	for _, a := range ids {
		if uint(a) == chatId {
			isAdmin = true
		}
	}

	if isAdmin && cmd.GetInput() == "/maintenance_on" && token.GetScenarioName() == "cryptoadmin" {
		setting.SetMaintenance(true)
		s.SettingsRepository.Persist(setting)
	}

	if isAdmin && cmd.GetInput() == "/maintenance_off" && token.GetScenarioName() == "cryptoadmin" {
		setting.SetMaintenance(false)
		s.SettingsRepository.Persist(setting)
	}

	return setting.IsMaintenance()
}
