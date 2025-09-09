package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (app *application) rateLimit() gin.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.Error(app.rateLimitExceededResponse("Too many requests, please try again in a few seconds"))
			return
		}
		c.Next()
	}
}
