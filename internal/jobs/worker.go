package jobs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"spotistats/internal/cache"
	"spotistats/internal/database"
)

func StartWorker(ctx context.Context) {
	log.Println("job worker started")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("job worker stopped")
				return
			default:
				processNextJob(context.Background())
			}
		}
	}()
}

func processNextJob(ctx context.Context) {
	// blocking pop from queue (5 second timeout)
	result, err := cache.Client.BRPop(ctx, 5*time.Second, JobQueueKey).Result()
	if err != nil {
		return // timeout or error, just retry
	}

	var job Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		log.Printf("[WORKER ERROR] unmarshal: %v", err)
		return
	}

	log.Printf("[JOB STARTED] id=%s type=%s", job.ID, job.Type)

	job.Status = JobStatusRunning
	updateJobStatus(ctx, &job)

	// process based on job type
	var jobErr error
	switch job.Type {
	case "genre_stats":
		job.Result, jobErr = processGenreStats(ctx)
	case "artist_stats":
		job.Result, jobErr = processArtistStats(ctx)
	case "year_stats":
		job.Result, jobErr = processYearStats(ctx)
	default:
		log.Printf("[WORKER ERROR] unknown job type: %s", job.Type)
		job.Status = JobStatusFailed
		updateJobStatus(ctx, &job)
		return
	}

	if jobErr != nil {
		log.Printf("[JOB FAILED] id=%s error=%v", job.ID, jobErr)
		job.Status = JobStatusFailed
		updateJobStatus(ctx, &job)
		return
	}

	job.Status = JobStatusCompleted
	updateJobStatus(ctx, &job)
	log.Printf("[JOB COMPLETED] id=%s type=%s", job.ID, job.Type)
}

func processGenreStats(ctx context.Context) (interface{}, error) {
	query := `
		SELECT genre, COUNT(*) as count 
		FROM tracks 
		WHERE genre != '' 
		GROUP BY genre 
		ORDER BY count DESC 
		LIMIT 10
	`

	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []map[string]interface{}{}
	for rows.Next() {
		var genre string
		var count int
		if err := rows.Scan(&genre, &count); err != nil {
			return nil, err
		}
		stats = append(stats, map[string]interface{}{
			"genre": genre,
			"count": count,
		})
	}

	return stats, nil
}

func processArtistStats(ctx context.Context) (interface{}, error) {
	query := `
		SELECT artist, COUNT(*) as count 
		FROM tracks 
		GROUP BY artist 
		ORDER BY count DESC 
		LIMIT 10
	`

	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []map[string]interface{}{}
	for rows.Next() {
		var artist string
		var count int
		if err := rows.Scan(&artist, &count); err != nil {
			return nil, err
		}
		stats = append(stats, map[string]interface{}{
			"artist": artist,
			"count":  count,
		})
	}

	return stats, nil
}

func processYearStats(ctx context.Context) (interface{}, error) {
	query := `
		SELECT year, COUNT(*) as count 
		FROM tracks 
		WHERE year > 0 
		GROUP BY year 
		ORDER BY year DESC 
		LIMIT 20
	`

	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []map[string]interface{}{}
	for rows.Next() {
		var year int
		var count int
		if err := rows.Scan(&year, &count); err != nil {
			return nil, err
		}
		stats = append(stats, map[string]interface{}{
			"year":  year,
			"count": count,
		})
	}

	return stats, nil
}
