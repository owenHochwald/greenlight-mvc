package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}

func newAppError(message string, code int) *AppError {
	return &AppError{message, code}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*AppError); ok {
				c.JSON(appErr.Code, gin.H{
					"success": false,
					"message": appErr.Message,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}
