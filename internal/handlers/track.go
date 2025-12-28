package handlers

import (
	"net/http"
	"strconv"
	"time"

	"spotistats/internal/models"
	"spotistats/internal/repository"

	"github.com/labstack/echo/v4"
)

// GetTracks godoc
// @Summary Get tracks
// @Description Get paginated list of tracks with optional filters
// @Tags tracks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param q query string false "Search by track name"
// @Param genre query string false "Filter by genre"
// @Param artist query string false "Filter by artist"
// @Success 200 {object} map[string]interface{}
// @Router /tracks [get]
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

// GetTrackByID godoc
// @Summary Get track by ID
// @Description Get a single track by track_id or spotify_id
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} models.Track
// @Failure 404 {object} map[string]string
// @Router /tracks/{id} [get]
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

// CreateTrack godoc
// @Summary Create track
// @Description Add a new track directly to database
// @Tags tracks
// @Accept json
// @Produce json
// @Param track body models.Track true "Track data"
// @Success 201 {object} models.Track
// @Failure 400 {object} map[string]string
// @Router /tracks [post]
func CreateTrack(c echo.Context) error {
	var track models.Track
	if err := c.Bind(&track); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if track.TrackID == "" {
		track.TrackID = "TRK" + time.Now().Format("20060102150405")
	}

	ctx := c.Request().Context()
	if err := repository.InsertTrack(ctx, track); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, track)
}
