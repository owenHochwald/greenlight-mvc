package main

import "C"
import (
	"errors"
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
		c.Error(app.badRequest("Invalid movie body"))
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
		c.Error(app.validationError("Make sure movie body is correct", v.Errors))
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
		c.Error(app.badRequest("Invalid movie parseID"))
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
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.NewValidator()

	qs := c.Request.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		c.Error(app.validationError("bad user input", v.Errors))
		c.Abort()
		return
	}

	movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)

	if err != nil {
		c.Error(newAppError("Database error", http.StatusInternalServerError, nil))
		c.Abort()
		return
	}

	c.IndentedJSON(200, gin.H{
		"movies":       movies,
		"current_page": input.Filters.Page,
		"page_size":    input.Filters.PageSize,
		"total_count":  len(movies),
	})
}

func (app *application) updateMovieHandler(c *gin.Context) {
	id, err := app.parseID(c)
	if err != nil {
		c.Error(app.badRequest("Invalid movie parseID"))
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
		c.Error(app.badRequest("Invalid movie data format"))
		c.Abort()
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		c.Error(app.databaseError("Failed to retrieve movie from database"))
		c.Abort()
		return
	}

	if movie == nil {
		c.Error(app.badRequest(fmt.Sprintf("Movie with ID %d not found", id)))
		c.Abort()
		return
	}

	applyMovieUpdates(movie, input)

	v := validator.NewValidator()
	data.ValidateMovie(v, movie)
	if !v.Valid() {
		c.Error(app.validationError("Movie validation failed", v.Errors))
		c.Abort()
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrMovieEditConflict):
			c.Error(app.editConflictError("Edit conflict: movie has been updated by another user"))
		default:
			c.Error(app.serverResponseError("Failed to update movie in database"))
		}
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"movie":   movie,
		"message": fmt.Sprintf("Successfully updated movie with ID %d", id),
	})
}

func (app *application) deleteMovieHandler(c *gin.Context) {
	id, err := app.parseID(c)

	if err != nil {
		c.Error(app.badRequest("Invalid movie parseID"))
		c.Abort()
		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		c.Error(app.databaseError(fmt.Sprintf("Failed to delete movie with parseID: %d", id)))
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
