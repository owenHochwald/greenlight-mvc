package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

const version = "1.0.0"

func main() {
	var cfg config

	// flag setup and parsing
	flag.IntVar(&cfg.port, "port", 8080, "port to run the server on")
	flag.StringVar(&cfg.env, "env", "dev", "environment to run the server in")
	flag.Parse()

	// setting up application and logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := application{
		config: cfg,
		logger: logger,
	}

	// creating gin app
	r := gin.Default()

	// attaching custom error handler
	r.Use(ErrorHandler())

	SetupRoutes(r, &app)

	r.Run(fmt.Sprintf(":%d", cfg.port))
}
