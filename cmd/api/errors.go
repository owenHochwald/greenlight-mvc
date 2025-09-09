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

func (app *application) badRequest(message string) *AppError {
	app.logger.Error().Msg("Bad user request")
	return newAppError(message, http.StatusBadRequest, nil)
}

func (app *application) databaseError(message string) *AppError {
	app.logger.Error().Msg("Database had unexpected error")
	return newAppError(message, http.StatusInternalServerError, nil)
}

func (app *application) validationError(message string, errors map[string]string) *AppError {
	app.logger.Error().Msg("Bad JSON object passed in request")
	return newAppError(message, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictError(message string) *AppError {
	app.logger.Warn().Msg("Database update edit conflict with concurrent requests")
	return newAppError(message, http.StatusConflict, nil)
}

func (app *application) serverResponseError(message string) *AppError {
	app.logger.Error().Msg("Server responded badly")
	return newAppError(message, http.StatusInternalServerError, nil)
}

func (app *application) rateLimitExceededResponse(message string) *AppError {
	app.logger.Error().Msg("Rate limit exceeded - Too many requests")
	return newAppError(message, http.StatusTooManyRequests, nil)
}

func newAppError(message string, code int, errors map[string]string) *AppError {
	return &AppError{
		message, code, errors,
	}
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
