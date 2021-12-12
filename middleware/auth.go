// Package middleware defines all the middlewares for the application
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const UserIDKey = "UserID"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(UserIDKey)
		if !exists {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
	}
}
