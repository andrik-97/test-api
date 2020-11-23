// Package task implements business use cases for task domain.
package task

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/payfazz/test-api/internal/task/model"
	"github.com/payfazz/test-api/internal/task/view"
	"github.com/payfazz/test-api/pkg/txn"
	"gopkg.in/go-playground/validator.v9"
)

// Service represents the interface for implementing task usecases.
type Service interface {
	CreateTask(context.Context, *CreateTaskInput) (*CreateTaskOutput, error)
	CompleteTask(context.Context, *CompleteTaskInput) error
	DeleteTask(context.Context, *DeleteTaskInput) error
	ViewTask(context.Context, *ViewTaskInput) (*ViewTaskOutput, error)
	ViewTasks(context.Context, *ViewTasksInput) (*ViewTasksOutput, error)
}

type service struct {
	taskRepo       model.TaskRepository
	taskView       view.TaskView
	validate       *validator.Validate
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	logger         log.Logger
	transactor     txn.Transactor
}

// ServiceConfig contains everything neede by the service.
type ServiceConfig struct {
	TaskRepo       model.TaskRepository
	TaskView       view.TaskView
	Validate       *validator.Validate
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Logger         log.Logger
	Transactor     txn.Transactor
}

// NewService creates a new task service.
func NewService(cfg *ServiceConfig) Service {
	return &service{
		taskRepo:       cfg.TaskRepo,
		taskView:       cfg.TaskView,
		validate:       cfg.Validate,
		requestCount:   cfg.RequestCount,
		requestLatency: cfg.RequestLatency,
		logger:         cfg.Logger,
		transactor:     cfg.Transactor,
	}
}
