package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/payfazz/backend-template/internal/task/model"
	"github.com/payfazz/backend-template/internal/task/view"
)

type (
	// CreateTaskInput represents the input parameters for creating a task.
	CreateTaskInput struct {
		Title string `json:"title"`
	}

	// CreateTaskOutput represents the output response for creating a task.
	CreateTaskOutput struct {
		TaskID string `json:"task_id"`
	}

	// CompleteTaskInput represents the input parameters for completing a task.
	CompleteTaskInput struct {
		TaskID string `httpurl:"task_id"`
	}

	// DeleteTaskInput represents the input parameters for deleting a task.
	DeleteTaskInput struct {
		TaskID string `httpurl:"task_id"`
	}

	// ViewTaskInput represents the input parameters for getting a task.
	ViewTaskInput struct {
		TaskID string `httpurl:"task_id"`
	}

	// ViewTaskOutput represents the output response for getting a task.
	ViewTaskOutput struct {
		*view.Task
	}

	// ViewTasksInput represents the parameters for listing all tasks.
	ViewTasksInput struct {
	}

	// ViewTasksOutput represents the output response for listing all tasks.
	ViewTasksOutput struct {
		Data []*view.Task `json:"data"`
	}
)

func (s *service) CreateTask(ctx context.Context, input *CreateTaskInput) (*CreateTaskOutput, error) {
	method := "create_task"
	defer func(begin time.Time) {
		s.requestCount.With("method", method).Add(1)
		s.requestLatency.With("method", method).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := s.validate.Struct(input)
	if err != nil {
		return nil, err
	}

	tsk := &model.Task{
		ID:    uuid.New().String(),
		Title: input.Title,
	}
	if err := s.taskRepo.Insert(ctx, tsk); err != nil {
		return nil, err
	}

	return &CreateTaskOutput{TaskID: tsk.ID}, nil
}

func (s *service) CompleteTask(ctx context.Context, input *CompleteTaskInput) error {
	method := "complete_task"
	defer func(begin time.Time) {
		s.requestCount.With("method", method).Add(1)
		s.requestLatency.With("method", method).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := s.validate.Struct(input)
	if err != nil {
		return err
	}

	err = s.transactor.RunInTransaction(ctx, func(txnCtx context.Context) error {
		tsk, err := s.taskRepo.Find(txnCtx, input.TaskID)
		if err != nil {
			return err
		}

		tsk.Complete()

		if err := s.taskRepo.Update(txnCtx, tsk); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *service) DeleteTask(ctx context.Context, input *DeleteTaskInput) error {
	method := "delete_task"
	defer func(begin time.Time) {
		s.requestCount.With("method", method).Add(1)
		s.requestLatency.With("method", method).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := s.validate.Struct(input)
	if err != nil {
		return err
	}

	err = s.transactor.RunInTransaction(ctx, func(txnCtx context.Context) error {
		tsk, err := s.taskRepo.Find(txnCtx, input.TaskID)
		if err != nil {
			return err
		}

		if err := s.taskRepo.Delete(txnCtx, tsk); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *service) ViewTask(ctx context.Context, input *ViewTaskInput) (*ViewTaskOutput, error) {
	method := "view_task"
	defer func(begin time.Time) {
		s.requestCount.With("method", method).Add(1)
		s.requestLatency.With("method", method).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := s.validate.Struct(input)
	if err != nil {
		return nil, err
	}

	tsk, err := s.taskView.Task(ctx, input.TaskID)
	return &ViewTaskOutput{Task: tsk}, err
}

func (s *service) ViewTasks(ctx context.Context, input *ViewTasksInput) (*ViewTasksOutput, error) {
	method := "view_tasks"
	defer func(begin time.Time) {
		s.requestCount.With("method", method).Add(1)
		s.requestLatency.With("method", method).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := s.validate.Struct(input)
	if err != nil {
		return nil, err
	}

	tsks, err := s.taskView.Tasks(ctx)
	return &ViewTasksOutput{Data: tsks}, err
}
