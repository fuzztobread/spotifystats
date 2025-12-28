package handlers

import (
	"net/http"
	"strconv"
	"spotistats/internal/kafka"
	"spotistats/internal/models"
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
func CreateTrack(c echo.Context) error {
	var track models.Track
	if err := c.Bind(&track); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// publish to kafka
	ctx := c.Request().Context()
	if err := kafka.Publish(ctx, track.TrackID, track); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to publish",
		})
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"status":   "queued",
		"track_id": track.TrackID,
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
