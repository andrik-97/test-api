// Package inmem implements domain model repositories in memory.
package inmem

import (
	"context"

	"github.com/payfazz/backend-template/internal/task/model"
)

// TaskRepository implements in-memory task.Repository.
type TaskRepository struct {
	Data map[string]*model.Task
}

// Find gets a task by its id from the inmem db.
func (r *TaskRepository) Find(ctx context.Context, id string) (*model.Task, error) {
	return r.Data[id], nil
}

// Insert inserts a task to the inmem db.
func (r *TaskRepository) Insert(ctx context.Context, tsk *model.Task) error {
	r.Data[tsk.ID] = tsk
	return nil
}

// Update updates a task within the inmem db.
func (r *TaskRepository) Update(ctx context.Context, tsk *model.Task) error {
	r.Data[tsk.ID] = tsk
	return nil
}

// Delete deletes a task from the inmem db.
func (r *TaskRepository) Delete(ctx context.Context, tsk *model.Task) error {
	delete(r.Data, tsk.ID)
	return nil
}
