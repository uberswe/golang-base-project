package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) NoRoute(c *gin.Context) {
	pd := PageData{
		Title:           "404 Not Found",
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "404.html", pd)
}
