package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Settings struct {
	gorm.Model
	ScenarioName                       string
	TelegramAdminsIds                  string
	TelegramNotificationChannelsTokens string
	CreatedAt                          time.Time // column name is `created_at`
	UpdatedAt                          time.Time // column name is `updated_at`
}

func (t *Settings) GetTelegramAdminsIds() []int {
	var objmap []int
	err := json.Unmarshal([]byte(t.TelegramAdminsIds), &objmap)

	if err != nil {
		panic(err)
	}

	return objmap
}

func (t *Settings) SetTelegramAdminsIds(ids []int) {
	result, err := json.Marshal(ids)

	if err != nil {
		panic(err)
	}

	t.TelegramAdminsIds = string(result)
}

func (t *Settings) GetTelegramNotificationChannelsTokens() []string {
	var objmap []string
	err := json.Unmarshal([]byte(t.TelegramNotificationChannelsTokens), &objmap)

	if err != nil {
		panic(err)
	}

	return objmap
}

func (t *Settings) SetTelegramNotificationChannelsTokens(ids []string) {
	result, err := json.Marshal(ids)

	if err != nil {
		panic(err)
	}

	t.TelegramNotificationChannelsTokens = string(result)
}

type SettingsRepository struct {
	DB *gorm.DB
}

func (r *SettingsRepository) Persist(settings *Settings) {
	r.DB.Save(settings)
}

func (r *SettingsRepository) Delete(settings *Settings) {
	r.DB.Delete(settings)
}

func (r *SettingsRepository) FindByScenarioName(scenarioName string) *Settings {
	var settings Settings
	res := r.DB.First(&settings, "scenario_name = ? and deleted_at IS NULL", scenarioName)
	if res != nil && res.Error != nil {
		fmt.Println(res.Error.Error())
		return nil
	}
	return &settings
}
