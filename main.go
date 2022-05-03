package main

import (
	"bot-daedalus/app"
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
)

type Config struct {
	DriverName           string
	DSN                  string
	PreferSimpleProtocol bool
	WithoutReturning     bool
	Conn                 gorm.ConnPool
}

func logRequestMiddleware(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}

func main() {
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

	if err != nil {
		panic(err)
	}

	healthcheckHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "service is working (ver. 1.2)",
		})
	}

	cryptoClientHandler := func(c *gin.Context) {
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
		logRequestMiddleware(c)
		bot.HandleRequest(&runtime.DefaultSerializedMessageFactory{Ctx: c})
	}

	cryptoAdminHandler := func(c *gin.Context) {
		actionRegistry := app.ActionRegistry{DB: db}
		commandRegistry := app.CommandRegistry{DB: db}
		bot := runtime.DefaultBot{
			ScenarioPath:       "config/scenario",
			ScenarioName:       "cryptoadmin",
			TokenFactory:       models.TokenFactory{DB: db},
			TokenRepository:    &models.TokenRepository{DB: db},
			SettingsRepository: &models.SettingsRepository{DB: db},
			ActionRegistry:     actionRegistry.ActionRegistryHandler,
			CommandRegistry:    commandRegistry.CommandRegistryHandler,
		}
		logRequestMiddleware(c)
		bot.HandleRequest(&runtime.DefaultSerializedMessageFactory{Ctx: c})
	}

	r := gin.Default()
	r.StaticFile("/1.jpg", "./resources/1.jpg")
	r.StaticFile("/2.jpg", "./resources/2.jpg")
	r.StaticFile("/3.jpg", "./resources/3.jpg")
	r.StaticFile("/4.jpg", "./resources/4.jpg")
	r.StaticFile("/5.jpg", "./resources/5.jpg")
	r.StaticFile("/6.jpg", "./resources/6.jpg")
	r.POST("/crypto", cryptoClientHandler)
	r.POST("/cryptoadmin", cryptoAdminHandler)
	r.GET("/", healthcheckHandler)

	err = r.Run(":8181")
	if err != nil {
		return
	}
}
