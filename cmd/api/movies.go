package main

import "C"
import (
	"fmt"
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

func (app *application) showAllMoviesHandler(c *gin.Context) {
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

func (app *application) updateMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		c.Error(badRequest("Invalid id passed"))
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err = c.ShouldBind(&input); err != nil {
		if err != nil {
			c.Error(badRequest("Invalid movie object passed"))
		}
	}

	movie, err := app.models.Movies.Get(id)

	if err != nil || movie == nil {
		c.Error(databaseError("Database error"))
		c.Abort()
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	data.ValidateMovie(validator.NewValidator(), movie)

	err = app.models.Movies.Update(movie)

	if err != nil {
		c.Error(databaseError("Database error"))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("successfully updated movie with id %s", movieId),
	})
}

func (app *application) deleteMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		c.Error(badRequest("Invalid id passed"))
		c.Abort()
		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		c.Error(databaseError(fmt.Sprintf("Failed to delete movie with id: %s", movieId)))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("successfully deleted movie with id %s", movieId),
	})

}
