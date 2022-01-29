package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NoRoute handles rendering of the 404 page
func (controller Controller) NoRoute(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("404 Not Found")
	c.HTML(http.StatusOK, "404.html", pd)
}
