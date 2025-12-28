package handlers

import (
	"net/http"

	"spotistats/internal/repository"

	"github.com/labstack/echo/v4"
)

func GetGenreStats(c echo.Context) error {
	stats, err := repository.GetGenreStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, stats)
}

func GetArtistStats(c echo.Context) error {
	stats, err := repository.GetArtistStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, stats)
}

func GetYearStats(c echo.Context) error {
	stats, err := repository.GetYearStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, stats)
}
