// Package postgres implements domain model repositories in postgresql db.
package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	pgmodel "github.com/payfazz/test-api/internal/postgres/model"
	taskmodel "github.com/payfazz/test-api/internal/task/model"
	"github.com/payfazz/test-api/pkg/txn/txnsql"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// TaskRepository implements task.Repository with postgres.
type TaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db}
}

// Find gets a task by its id from the postgres db.
func (r *TaskRepository) Find(ctx context.Context, id string) (*taskmodel.Task, error) {
	pgTask, err := pgmodel.FindTask(ctx, r.db, id)
	if err != nil {
		return nil, err
	}

	return encodeToTask(pgTask), nil
}

// Insert inserts a task to the postgres db.
func (r *TaskRepository) Insert(ctx context.Context, tsk *taskmodel.Task) error {
	db := boil.ContextExecutor(r.db)
	if tx, ok := txnsql.TxFromContext(ctx); ok {
		db = tx
	}

	pgTask := decodeTask(tsk)

	if err := pgTask.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

// Update updates a task within the postgres db.
func (r *TaskRepository) Update(ctx context.Context, tsk *taskmodel.Task) error {
	db := boil.ContextExecutor(r.db)
	if tx, ok := txnsql.TxFromContext(ctx); ok {
		db = tx
	}

	pgTask := decodeTask(tsk)

	if _, err := pgTask.Update(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

// Delete deletes a task from the postgres db.
func (r *TaskRepository) Delete(ctx context.Context, tsk *taskmodel.Task) error {
	db := boil.ContextExecutor(r.db)
	if tx, ok := txnsql.TxFromContext(ctx); ok {
		db = tx
	}

	pgTask, err := pgmodel.FindTask(ctx, db, tsk.ID)
	if err != nil {
		return err
	}

	if _, err := pgTask.Delete(ctx, db, false); err != nil {
		return err
	}

	return nil
}

func decodeTask(tsk *taskmodel.Task) *pgmodel.Task {
	return &pgmodel.Task{
		ID:          tsk.ID,
		Title:       tsk.Title,
		Completed:   tsk.Completed,
		CompletedAt: null.TimeFromPtr(tsk.CompletedAt),
	}
}

func encodeToTask(pgTask *pgmodel.Task) *taskmodel.Task {
	return &taskmodel.Task{
		ID:          pgTask.ID,
		Title:       pgTask.Title,
		Completed:   pgTask.Completed,
		CompletedAt: pgTask.CompletedAt.Ptr(),
	}
}
