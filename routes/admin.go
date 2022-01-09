package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Admin renders the admin dashboard
func (controller Controller) Admin(c *gin.Context) {
	pd := PageData{
		Title:           "Admin",
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  controller.config.CacheParameter,
	}
	c.HTML(http.StatusOK, "admin.html", pd)
}
