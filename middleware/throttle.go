package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	middlewareGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"log"
)

func Throttle(limit int) gin.HandlerFunc {
	// Define a limit rate to 3 requests per minute.
	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		log.Fatal(err)
	}

	store := memory.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new middleware with the limiter instance.
	return middlewareGin.NewMiddleware(limiter.New(store, rate))
}
