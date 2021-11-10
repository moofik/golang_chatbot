package main

import (
	"bot-daedalus/app"
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Token{})
	err = db.AutoMigrate(&models.Order{})
	err = db.AutoMigrate(&models.Wallet{})
	err = db.AutoMigrate(&models.WalletOrder{})

	if err != nil {
		panic(err)
	}

	handler := func(c *gin.Context) {
		//fmt.Printf(newStr)
		//provider.GetCommandFromRequest(c)
		actionRegistry := app.ActionRegistry{DB: db}
		bot := runtime.DefaultBot{
			ScenarioPath:    "config/scenario",
			ScenarioName:    "scenario",
			TokenFactory:    models.TokenFactory{DB: db},
			TokenRepository: &models.TokenRepository{DB: db},
			ActionRegistry:  actionRegistry.ActionRegistryHandler,
		}

		//fmt.Println("REQUEST BODY:")
		//bot.LogRequest(c)

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
