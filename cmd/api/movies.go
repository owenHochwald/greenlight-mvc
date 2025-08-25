package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"owenHochwald.greenlight/internal/data"
)

func (app *application) createMovieHandler(c *gin.Context) {
	c.IndentedJSON(200, gin.H{
		"message": "Hello World",
	})

}

func (app *application) showMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		c.IndentedJSON(400, gin.H{
			"message": "Invalid movie id",
		})
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
