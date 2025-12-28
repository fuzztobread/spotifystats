package commands

import (
	"log"
	"net/http"

	"spotistats/internal/cache"
	"spotistats/internal/config"
	"spotistats/internal/database"
	"spotistats/internal/handlers"

	_ "spotistats/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// stats routes
	e.GET("/stats/genres", handlers.GetGenreStats)
	e.GET("/stats/artists", handlers.GetArtistStats)
	e.GET("/stats/years", handlers.GetYearStats)

	log.Printf("starting server on :%s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
