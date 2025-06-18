package domain

import (
	"errors"
	"time"
)

type TaskFunction func() (interface{}, error)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID          string
	Status      TaskStatus
	Result      interface{}
	Error       error
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
}

func NewTask(id string) *Task {
	return &Task{
		ID:        id,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

func (t *Task) Start() {
	now := time.Now()
	t.StartedAt = &now
	t.Status = StatusInProgress
}

func (t *Task) Complete(result interface{}, err error) {
	now := time.Now()
	t.CompletedAt = &now
	t.Result = result
	t.Error = err

	if err != nil {
		t.Status = StatusFailed
	} else {
		t.Status = StatusCompleted
	}
}

func (t *Task) Duration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}

	if t.CompletedAt != nil {
		return t.CompletedAt.Sub(*t.StartedAt)
	}

	return time.Since(*t.StartedAt)
}
