package main

import (
	"bot-daedalus/app"
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

func cli() {
	err := godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Token{})
	err = db.AutoMigrate(&models.Order{})
	err = db.AutoMigrate(&models.Wallet{})
	err = db.AutoMigrate(&models.WalletOrder{})
	err = db.AutoMigrate(&models.Settings{})

	actionRegistry := app.ActionRegistry{DB: db}
	commandRegistry := app.CommandRegistry{DB: db}
	bot := runtime.DefaultBot{
		ScenarioPath:       "config/scenario",
		ScenarioName:       "cryptobot",
		TokenFactory:       models.TokenFactory{DB: db},
		TokenRepository:    &models.TokenRepository{DB: db},
		SettingsRepository: &models.SettingsRepository{DB: db},
		ActionRegistry:     actionRegistry.ActionRegistryHandler,
		CommandRegistry:    commandRegistry.CommandRegistryHandler,
		StateErrorHandler:  app.CryptobotStateErrorHandler,
	}

	cmd := os.Args[1]

	if cmd == "run" {
		argActionName := os.Args[2]
		argTokenId := os.Args[3]

		tokenRepository := &models.TokenRepository{DB: db}
		i, _ := strconv.ParseInt(argTokenId, 10, 64)
		token := tokenRepository.FindById(int(i))
		_, provider, scenario := bot.GetBaseActors(&runtime.DefaultSerializedMessageFactory{Ctx: nil})
		currentState := scenario.GetCurrentState(token)

		a1 := app.CalculateMarketBuyOrder{
			OrderRepository: &models.OrderRepository{DB: db},
		}

		if argActionName == a1.GetName() || argActionName == a1.GetAlias() {
			a1.Run(provider, token, currentState, nil, nil)
			tokenRepository.Persist(token)
		}
	}

	if cmd == "next" {
		argTokenId := os.Args[2]

		tokenRepository := &models.TokenRepository{DB: db}
		i, _ := strconv.ParseInt(argTokenId, 10, 64)
		token := tokenRepository.FindById(int(i))
		_, _, scenario := bot.GetBaseActors(&runtime.DefaultSerializedMessageFactory{Ctx: nil})
		currentState := scenario.GetCurrentState(token)

		cmd, _ := currentState.GetCommandByProto(&runtime.ButtonPressedCommand{
			ButtonCommand: "",
			ButtonText:    "",
			Metadata: &runtime.CommandMetadata{
				Cmd: "button",
			},
		})

		scenario.HandleCommand(
			cmd,
			currentState,
			token,
		)
	}

}
