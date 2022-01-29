package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Admin renders the admin dashboard
func (controller Controller) Admin(c *gin.Context) {
	pd := controller.DefaultPageData(c)
	pd.Title = pd.Trans("Admin")
	c.HTML(http.StatusOK, "admin.html", pd)
}
