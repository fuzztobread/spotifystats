package main

import (
	"context"
	"log"
	"net/http"

	"spotistats/internal/cache"
	"spotistats/internal/database"
	"spotistats/internal/handlers"
	"spotistats/internal/jobs"
	"spotistats/internal/kafka"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// connect to postgres
	err := database.Connect("postgres://spotistats:spotistats@localhost:5432/spotistats")
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer database.Close()

	// connect to redis
	err = cache.Connect("localhost:6379")
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	defer cache.Close()

	// init kafka
	brokers := []string{"localhost:9092"}
	kafka.InitProducer(brokers, "tracks_ingest")
	defer kafka.CloseProducer()

	// start kafka consumer
	kafka.StartConsumer(brokers, "tracks_ingest", "spotistats-group")

	// start job worker
	jobs.StartWorker(context.Background())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// track routes
	e.GET("/tracks", handlers.GetTracks)
	e.GET("/tracks/:id", handlers.GetTrackByID)
	e.POST("/tracks", handlers.CreateTrack)

	// job routes
	e.POST("/jobs", handlers.CreateJob)
	e.GET("/jobs/:id", handlers.GetJobStatus)

	e.Logger.Fatal(e.Start(":8080"))
}
