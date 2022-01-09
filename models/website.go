package models

import "gorm.io/gorm"

// Website holds information about different websites
type Website struct {
	gorm.Model
	Title       string
	Description string
	URL         string
}
