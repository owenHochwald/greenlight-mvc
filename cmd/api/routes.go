package main

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, app *application) {

	r.GET("/v1/healthcheck", app.healthCheckHandler)
	r.GET("/v1/movies/:id", app.showMovieHandler)
	r.POST("/v1/movies", app.createMovieHandler)

}
