// Package model provides domain models in task domain.
package model

import (
	"context"
	"time"
)

// Task represents a todo.
type Task struct {
	ID          string
	Title       string
	Completed   bool
	CompletedAt *time.Time
}

// Complete marks the task as completed.
func (t *Task) Complete() {
	currTime := time.Now()
	t.Completed = true
	t.CompletedAt = &currTime
}

// TaskRepository represents an interface for task storage operations.
type TaskRepository interface {
	Find(ctx context.Context, id string) (*Task, error)
	Insert(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, task *Task) error
}
