package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) Admin(c *gin.Context) {
	pd := PageData{
		Title:           "Admin",
		IsAuthenticated: isAuthenticated(c),
	}
	c.HTML(http.StatusOK, "admin.html", pd)
}
