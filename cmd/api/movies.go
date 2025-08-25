package main

import "github.com/gin-gonic/gin"

func (app *application) createMovieHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})

}

func (app *application) showMovieHandler(c *gin.Context) {
	movieId := c.Param("id")
	c.JSON(200, gin.H{
		"message": "Hello World",
		"id":      movieId,
	})
}
