package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Cache(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Next()
	}
}
