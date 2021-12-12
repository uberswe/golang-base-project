package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Identifier string
	UserID     uint
}
