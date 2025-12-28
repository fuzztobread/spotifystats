package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"spotistats/internal/models"

	"github.com/labstack/echo/v4"
)

var Tracks []models.Track

func GetTracks(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// search params
	query := strings.ToLower(c.QueryParam("q"))
	genre := strings.ToLower(c.QueryParam("genre"))
	artist := strings.ToLower(c.QueryParam("artist"))

	// filter tracks
	var filtered []models.Track
	for _, t := range Tracks {
		if query != "" && !strings.Contains(strings.ToLower(t.Name), query) {
			continue
		}
		if genre != "" && !strings.Contains(strings.ToLower(t.Genre), genre) {
			continue
		}
		if artist != "" && !strings.Contains(strings.ToLower(t.Artist), artist) {
			continue
		}
		filtered = append(filtered, t)
	}

	// paginate
	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total":   total,
		"page":    page,
		"limit":   limit,
		"results": filtered[start:end],
	})
}

func GetTrackByID(c echo.Context) error {
	id := c.Param("id")

	for _, track := range Tracks {
		if track.TrackID == id || track.SpotifyID == id {
			return c.JSON(http.StatusOK, track)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "track not found",
	})
}
