package main

import (
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
	if err != nil {
		return
	}

	handler := func(c *gin.Context) {
		//fmt.Printf(newStr)
		//provider.GetCommandFromRequest(c)

		bot := runtime.DefaultBot{
			ScenarioPath:    "config/scenario",
			ScenarioName:    "scenario",
			TokenFactory:    models.TokenFactory{DB: db},
			TokenRepository: models.TokenRepository{DB: db},
		}

		//fmt.Println("REQUEST BODY:")
		//bot.LogRequest(c)

		bot.HandleRequest(&runtime.DefaultSerializedMessageFactory{Ctx: c})

	}

	r := gin.Default()
	r.POST("/", handler)

	err = r.Run()
	if err != nil {
		return
	}
}
