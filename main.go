package main

import (
	"bot-daedalus/app"
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	DriverName           string
	DSN                  string
	PreferSimpleProtocol bool
	WithoutReturning     bool
	Conn                 gorm.ConnPool
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

	handler := func(c *gin.Context) {
		actionRegistry := app.ActionRegistry{DB: db}
		bot := runtime.DefaultBot{
			ScenarioPath:    "config/scenario",
			ScenarioName:    "cryptobot",
			TokenFactory:    models.TokenFactory{DB: db},
			TokenRepository: &models.TokenRepository{DB: db},
			ActionRegistry:  actionRegistry.ActionRegistryHandler,
		}

		bot.HandleRequest(&runtime.DefaultSerializedMessageFactory{Ctx: c})
	}

	healthcheckHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "service is working",
		})
	}

	r := gin.Default()
	r.POST("/crypto", handler)
	r.GET("/", healthcheckHandler)

	err = r.Run(":8181")
	if err != nil {
		return
	}
}
