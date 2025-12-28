package main

import (
	"log"
	"net/http"

	"spotistats/internal/handlers"
	"spotistats/internal/loader"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	tracks, err := loader.LoadTracksFromCSV("data/tracks.csv")
	if err != nil {
		log.Fatal("failed to load CSV:", err)
	}
	log.Printf("loaded %d tracks", len(tracks))

	handlers.Tracks = tracks

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.GET("/tracks", handlers.GetTracks)
	e.GET("/tracks/:id", handlers.GetTrackByID)

	e.Logger.Fatal(e.Start(":8080"))
}
