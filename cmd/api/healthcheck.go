package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) healthCheckHandler(c *gin.Context) {
	c.IndentedJSON(200, gin.H{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	})
}
