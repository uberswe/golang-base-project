package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	middlewareGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

func Throttle(limit int) gin.HandlerFunc {
	store := memory.NewStore()

	// Create a new middleware with the limiter instance.
	return middlewareGin.NewMiddleware(limiter.New(store, limiter.Rate{
		Period: time.Minute,
		Limit:  int64(limit),
	}))
}
