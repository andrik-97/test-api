// Package view provides read-only models for displaying the data to users.
package view

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	pgmodel "github.com/payfazz/test-api/internal/postgres/model"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Task is a default read-only model for task.
	Task struct {
		ID          string     `json:"id"`
		Title       string     `json:"title"`
		Completed   bool       `json:"completed"`
		CompletedAt *time.Time `json:"completed_at"`
	}
)

// TaskView represents view operations for task.
type TaskView interface {
	Task(ctx context.Context, taskID string) (*Task, error)
	Tasks(ctx context.Context) ([]*Task, error)
}

type taskView struct {
	db       *sqlx.DB
	validate *validator.Validate
}

// NewTaskView creates a new task view.
func NewTaskView(db *sqlx.DB, validate *validator.Validate) TaskView {
	return &taskView{db, validate}
}

// Task gets a task.
func (v *taskView) Task(ctx context.Context, taskID string) (*Task, error) {
	pgTask, err := pgmodel.FindTask(ctx, v.db, taskID)
	if err != nil {
		return nil, err
	}

	return encodeToTask(pgTask), nil
}

// Tasks lists tasks.
func (v *taskView) Tasks(ctx context.Context) ([]*Task, error) {
	pgTasks, err := pgmodel.Tasks().All(ctx, v.db)
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	for _, pgTask := range pgTasks {
		tasks = append(tasks, encodeToTask(pgTask))
	}

	return tasks, nil
}

func encodeToTask(pgTask *pgmodel.Task) *Task {
	return &Task{
		ID:          pgTask.ID,
		Title:       pgTask.Title,
		Completed:   pgTask.Completed,
		CompletedAt: pgTask.CompletedAt.Ptr(),
	}
}
