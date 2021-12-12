package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Value     string
	Type      string
	ModelID   int
	ModelType string
}

const (
	TokenUserActivation string = "user_activation"
	TokenPasswordReset  string = "password_reset"
)
