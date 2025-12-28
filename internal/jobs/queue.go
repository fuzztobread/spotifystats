package jobs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"spotistats/internal/cache"
)

type Job struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Status    string                 `json:"status"`
	Result    interface{}            `json:"result,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"

	JobQueueKey = "jobs:queue"
)

func Enqueue(ctx context.Context, job Job) error {
	job.Status = JobStatusPending
	job.CreatedAt = time.Now()

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// save job status
	if err := cache.Client.Set(ctx, "job:"+job.ID, data, 1*time.Hour).Err(); err != nil {
		return err
	}

	// push to queue
	if err := cache.Client.LPush(ctx, JobQueueKey, data).Err(); err != nil {
		return err
	}

	log.Printf("[JOB QUEUED] id=%s type=%s", job.ID, job.Type)
	return nil
}

func GetJobStatus(ctx context.Context, jobID string) (*Job, error) {
	data, err := cache.Client.Get(ctx, "job:"+jobID).Bytes()
	if err != nil {
		return nil, err
	}

	var job Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func updateJobStatus(ctx context.Context, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return cache.Client.Set(ctx, "job:"+job.ID, data, 1*time.Hour).Err()
}
