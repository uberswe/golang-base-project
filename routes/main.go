// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/uberswe/golang-base-project/config"
	"github.com/uberswe/golang-base-project/lang"
	"github.com/uberswe/golang-base-project/middleware"
	"gorm.io/gorm"
)

// Controller holds all the variables needed for routes to perform their logic
type Controller struct {
	db     *gorm.DB
	config config.Config
	bundle *i18n.Bundle
}

// New creates a new instance of the routes.Controller
func New(db *gorm.DB, c config.Config, bundle *i18n.Bundle) Controller {
	return Controller{
		db:     db,
		config: c,
		bundle: bundle,
	}
}

// PageData holds the default data needed for HTML pages to render
type PageData struct {
	Title           string
	Messages        []Message
	IsAuthenticated bool
	CacheParameter  string
	Trans           func(s string) string
}

// Message holds a message which can be rendered as responses on HTML pages
type Message struct {
	Type    string // success, warning, error, etc.
	Content string
}

// isAuthenticated checks if the current user is authenticated or not
func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	return exists
}

func (controller Controller) DefaultPageData(c *gin.Context) PageData {
	langService := lang.New(c, controller.bundle)
	return PageData{
		Title:           "Home",
		Messages:        nil,
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
		Trans:           langService.Trans,
	}
}
