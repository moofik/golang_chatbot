package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	ChatId uint
	State string
	ScenarioName string
	IsBlocked bool
	TimeOffset int
	UserName string
	FirstName string
	LastName string
}
