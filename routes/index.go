package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index renders the HTML of the index page
func (controller Controller) Index(c *gin.Context) {
	pd := PageData{
		Title:           "Home",
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "index.html", pd)
}
