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
	db.AutoMigrate(&models.Token{})

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {

		//fmt.Printf(newStr)
		//provider.GetCommandFromRequest(c)

		bot := runtime.DefaultBot{}
		bot.HandleRequest(c)

		c.JSON(200, gin.H{
			"message": "hello",
		})
	})

	r.Run()
}
