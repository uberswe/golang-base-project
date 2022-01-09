package models

import (
	"gorm.io/gorm"
	"time"
)

// Session holds information about user sessions and when they expire
type Session struct {
	gorm.Model
	Identifier string
	UserID     uint
	ExpiresAt  time.Time
}

// HasExpired is a helper function that checks if the current time is after the session expire datetime
func (s Session) HasExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
