package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NoAuth is for routes that can only be accessed when the user is unauthenticated
func NoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(UserIDKey)
		if !exists {
			c.Next()
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "/admin")
		c.Abort()
	}
}
