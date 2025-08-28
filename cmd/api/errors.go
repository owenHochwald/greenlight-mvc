package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func badRequest(message string) *AppError {
	return newAppError(message, http.StatusBadRequest, nil)
}

func databaseError(message string) *AppError {
	return newAppError(message, http.StatusInternalServerError, nil)
}

func newAppError(message string, code int, errors map[string]string) *AppError {
	return &AppError{
		message, code, errors,
	}
}

func validationError(message string, errors map[string]string) *AppError {
	return newAppError(message, http.StatusUnprocessableEntity, errors)
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*AppError); ok {
				response := gin.H{
					"success": false,
					"message": appErr.Message,
				}
				if appErr.Errors != nil && len(appErr.Errors) > 0 {
					response["errors"] = appErr.Errors
				}
				c.JSON(appErr.Code, response)
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}
