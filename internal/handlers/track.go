package handlers

import (
	"net/http"
	"strconv"

	"spotistats/internal/repository"

	"github.com/labstack/echo/v4"
)

func GetTracks(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := c.QueryParam("q")
	genre := c.QueryParam("genre")
	artist := c.QueryParam("artist")

	offset := (page - 1) * limit

	tracks, total, err := repository.GetTracks(c.Request().Context(), limit, offset, query, genre, artist)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total":   total,
		"page":    page,
		"limit":   limit,
		"results": tracks,
	})
}

func GetTrackByID(c echo.Context) error {
	id := c.Param("id")

	track, err := repository.GetTrackByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "track not found",
		})
	}

	return c.JSON(http.StatusOK, track)
}
