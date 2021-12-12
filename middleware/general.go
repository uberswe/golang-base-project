package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

// General handles the default headers that should be present in every response
func General() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Frame-Options", "DENY")
		if os.Getenv("STRICT_TRANSPORT_SECURITY") == "true" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload;")
		}
		c.Next()
	}
}
