package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Cache middleware sets the Cache-Control header
func Cache(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Next()
	}
}
