package handlers

import (
	"net/http"

	"spotistats/internal/jobs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateJobRequest struct {
	Type string `json:"type"`
}

func CreateJob(c echo.Context) error {
	var req CreateJobRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	// validate job type
	validTypes := map[string]bool{
		"genre_stats":  true,
		"artist_stats": true,
		"year_stats":   true,
	}

	if !validTypes[req.Type] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid job type",
		})
	}

	job := jobs.Job{
		ID:   uuid.New().String(),
		Type: req.Type,
	}

	ctx := c.Request().Context()
	if err := jobs.Enqueue(ctx, job); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to queue job",
		})
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"job_id": job.ID,
		"status": "queued",
	})
}

func GetJobStatus(c echo.Context) error {
	jobID := c.Param("id")

	job, err := jobs.GetJobStatus(c.Request().Context(), jobID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "job not found",
		})
	}

	return c.JSON(http.StatusOK, job)
}
