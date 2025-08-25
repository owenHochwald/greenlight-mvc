package main

import "C"
import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"owenHochwald.greenlight/internal/data"
)

func (app *application) createMovieHandler(c *gin.Context) {

	var input struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := app.readJSON(c, &input); err != nil {
		c.Error(badRequest("Invalid movie body"))
	}

	//if err := c.ShouldBind(&input); err != nil {
	//	c.Error(badRequest("Invalid movie body"))
	//}

	c.IndentedJSON(http.StatusOK, gin.H{
		"movie":   input,
		"message": "Successfully created movie",
	})
}

func (app *application) showMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		c.Error(newAppError("Invalid movie id", 400))
	}

	movie := data.Movie{
		ID:       id,
		CreateAt: time.Now(),
		Title:    "test",
		Year:     2025,
		Runtime:  120,
		Genres:   []string{"test"},
		Version:  1,
	}

	c.IndentedJSON(200, gin.H{
		"movie": movie,
	})
}
