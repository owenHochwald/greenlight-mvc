package main

import "C"
import (
	"fmt"
	"net/http"

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
	id, err := app.parseID(c)

	if err != nil {
		c.Error(badRequest("Invalid movie parseID"))
		c.Abort()
		return
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
	id, err := app.parseID(c)

	if err != nil {
		c.Error(badRequest("Invalid movie parseID"))
		c.Abort()
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	if err = app.readJSON(c, &input); err != nil {
		if err != nil {
			c.Error(badRequest("Invalid movie object passed"))
		}
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		if movie == nil {
			c.Error(badRequest(err.Error()))
		} else {
			c.Error(databaseError("Failed to retrieve movie from database"))
		}
		c.Abort()
		return
	}

	applyMovieUpdates(movie, input)

	data.ValidateMovie(validator.NewValidator(), movie)

	err = app.models.Movies.Update(movie)

	if err != nil {
		c.Error(databaseError("Database error"))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("successfully updated movie with parseID %d", id),
	})
}

func (app *application) deleteMovieHandler(c *gin.Context) {
	id, err := app.parseID(c)

	if err != nil {
		c.Error(badRequest("Invalid movie parseID"))
		c.Abort()
		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		c.Error(databaseError(fmt.Sprintf("Failed to delete movie with parseID: %d", id)))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("successfully deleted movie with parseID %d", id),
	})
}

func applyMovieUpdates(movie *data.Movie, input struct {
	Title   *string       `json:"title"`
	Year    *int32        `json:"year"`
	Runtime *data.Runtime `json:"runtime"`
	Genres  []string      `json:"genres"`
}) {
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}
}
