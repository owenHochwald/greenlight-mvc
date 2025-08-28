package main

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, app *application) {

	r.GET("/v1/healthcheck", app.healthCheckHandler)
	r.GET("/v1/movies/:id", app.showMovieHandler)
	r.POST("/v1/movies", app.createMovieHandler)
	r.GET("/v1/movies", app.showAllMoviesHandler)
	r.PUT("/v1/movies/:id", app.updateMovieHandler)
	r.DELETE("/v1/movies/:id", app.deleteMovieHandler)
}
