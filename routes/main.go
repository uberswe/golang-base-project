// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/config"
	"github.com/uberswe/golang-base-project/middleware"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	config config.Config
}

func New(db *gorm.DB, c config.Config) Controller {
	return Controller{
		db:     db,
		config: c,
	}
}

type PageData struct {
	Title           string
	Messages        []Message
	IsAuthenticated bool
	CacheParameter  string
}

type Message struct {
	Type    string // success, warning, error, etc.
	Content string
}

func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	return exists
}
