package middleware

import "github.com/gin-gonic/gin"

func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=60")
		c.Next()
	}
}
