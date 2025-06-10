package application

import "time"

type TaskDTO struct {
	ID          string        `json:"id"`
	Status      string        `json:"status"`
	Result      interface{}   `json:"result,omitempty"`
	Error       string        `json:"error,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
}
