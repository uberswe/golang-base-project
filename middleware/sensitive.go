package middleware

import (
	"github.com/gin-gonic/gin"
)

// Sensitive middleware handles headers that should be set for routes that may contain sensitive data such as personal information
func Sensitive() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Next()
	}
}
