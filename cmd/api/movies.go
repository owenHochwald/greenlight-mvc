package main

import "C"
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"owenHochwald.greenlight/internal/data"
	"owenHochwald.greenlight/internal/validator"
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
		c.Abort()
		return
	}

	v := validator.NewValidator()

	movie := &data.Movie{
		Title:   input.Title,
		Year:    int32(input.Year),
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		c.Error(validationError("Make sure movie body is correct", v.Errors))
		c.Abort()
		return
	}

	err := app.models.Movies.Insert(movie)

	if err != nil {
		c.Error(newAppError("Database error", http.StatusInternalServerError, nil))
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"movie":   input,
		"message": "Successfully created movie",
	})
}

func (app *application) showMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		c.Error(validationError("Invalid movie id", nil))
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil {
		c.Error(newAppError("Database error", http.StatusInternalServerError, nil))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"movie": movie,
	})
}

func (app *application) showALlMoviesHandler(c *gin.Context) {
	movies, err := app.models.Movies.GetAll()

	if err != nil {
		c.Error(newAppError("Database error", http.StatusInternalServerError, nil))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"movies": movies,
	})

}
