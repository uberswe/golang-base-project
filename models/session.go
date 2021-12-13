package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Identifier string
	UserID     uint
	ExpiresAt  time.Time
}

func (s Session) HasExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
