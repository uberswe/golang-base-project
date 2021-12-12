package models

import "gorm.io/gorm"

type Website struct {
	gorm.Model
	Title       string
	Description string
	URL         string
}
