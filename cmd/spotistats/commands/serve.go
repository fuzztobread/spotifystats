package commands

import (
	"context"
	"log"
	"net/http"

	"spotistats/internal/cache"
	"spotistats/internal/config"
	"spotistats/internal/database"
	"spotistats/internal/handlers"
	"spotistats/internal/jobs"
	"spotistats/internal/kafka"

	_ "spotistats/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	// connect to postgres
	if err := database.Connect(cfg.DatabaseURL()); err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer database.Close()

	// connect to redis
	if err := cache.Connect(cfg.RedisAddr); err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	defer cache.Close()

	// init kafka
	kafka.InitProducer(cfg.KafkaBrokers, "tracks_ingest")
	defer kafka.CloseProducer()

	// start kafka consumer
	kafka.StartConsumer(cfg.KafkaBrokers, "tracks_ingest", "spotistats-group")

	// start job worker
	jobs.StartWorker(context.Background())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// dashboard
	e.GET("/", handlers.ServeDashboard)
	e.Static("/static", "web/static")

	// health check
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

	log.Printf("starting server on :%s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
