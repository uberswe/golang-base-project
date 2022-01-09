package models

import (
	"gorm.io/gorm"
	"time"
)

// Token holds tokens typically used for user activation and password resets
type Token struct {
	gorm.Model
	Value     string
	Type      string
	ModelID   int
	ModelType string
	ExpiresAt time.Time
}

// HasExpired is a helper function that checks if the current time is after the token expiration time
func (t Token) HasExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

const (
	// TokenUserActivation is a constant used to identify tokens used for user activation
	TokenUserActivation string = "user_activation"
	// TokenPasswordReset is a constant used to identify tokens used for password resets
	TokenPasswordReset string = "password_reset"
)
