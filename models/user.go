// Package models defines all the database models for the application
package models

import (
	"gorm.io/gorm"
	"time"
)

// User holds information relating to users that use the application
type User struct {
	gorm.Model
	Email       string
	Password    string
	ActivatedAt *time.Time
	Tokens      []Token `gorm:"polymorphic:Model;"`
	Sessions    []Session
}
