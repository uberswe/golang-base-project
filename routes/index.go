package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller Controller) Index(c *gin.Context) {
	pd := PageData{
		Title:           "Home",
		IsAuthenticated: isAuthenticated(c),
	}
	c.HTML(http.StatusOK, "index.html", pd)
}
