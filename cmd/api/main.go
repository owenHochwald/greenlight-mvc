package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	//"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"owenHochwald.greenlight/internal/data"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger zerolog.Logger
	models *data.Models
}

const version = "1.0.0"

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}
	var cfg config

	// flag setup and parsing
	flag.IntVar(&cfg.port, "port", 8080, "port to run the server on")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENVIRONMENT"), "environment to run the server in")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DB_STRING"), "database connection string")
	flag.Parse()

	// connecting to database
	db, err := ConnectDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error connecting to database")
	}
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(10 * time.Minute)

	logger.Info().Msgf("Connected to Postgres")

	defer db.Close()

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	r := gin.Default()

	r.Use(app.rateLimit())

	SetupRoutes(r, &app)

	r.Run(fmt.Sprintf(":%d", cfg.port))
}

func ConnectDB(cfg config) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
