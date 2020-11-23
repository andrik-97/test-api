package task

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/payfazz/test-api/internal/task/model"
	"github.com/payfazz/test-api/internal/task/model/inmem"
)

func TestSuccessCreateTask(t *testing.T) {
	tests := map[string]struct {
		input *CreateTaskInput
		want  *model.Task
	}{
		"success": {
			input: &CreateTaskInput{Title: "Todo 1"},
			want: &model.Task{
				Title:       "Todo 1",
				Completed:   false,
				CompletedAt: nil,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			repo := &inmem.TaskRepository{Data: map[string]*model.Task{}}

			svc := testNewService(
				&ServiceConfig{TaskRepo: repo},
			)

			out, err := svc.CreateTask(context.TODO(), tt.input)
			if err != nil {
				t.Fatal(err)
			}

			tsk := repo.Data[out.TaskID]
			c.Assert(tsk, qt.CmpEquals(cmpopts.IgnoreFields(model.Task{}, "ID")), tt.want)
		})
	}
}

func TestSuccessCompleteTask(t *testing.T) {
	tests := map[string]struct {
		input *CompleteTaskInput
		tsk   *model.Task
		want  *model.Task
	}{
		"success": {
			input: &CompleteTaskInput{TaskID: "1"},
			tsk:   &model.Task{ID: "1", Completed: false},
			want:  &model.Task{ID: "1", Completed: true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			repo := &inmem.TaskRepository{Data: map[string]*model.Task{
				tt.tsk.ID: tt.tsk,
			}}

			svc := testNewService(
				&ServiceConfig{TaskRepo: repo},
			)

			err := svc.CompleteTask(context.TODO(), tt.input)
			if err != nil {
				t.Fatal(err)
			}

			tsk := repo.Data[tt.tsk.ID]
			c.Assert(tsk, qt.CmpEquals(cmpopts.IgnoreFields(model.Task{}, "CompletedAt")), tt.want)
		})
	}
}

func TestSuccessDeleteTask(t *testing.T) {
	tests := map[string]struct {
		input *DeleteTaskInput
		tsk   *model.Task
	}{
		"success": {
			input: &DeleteTaskInput{TaskID: "1"},
			tsk:   &model.Task{ID: "1", Completed: false},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			repo := &inmem.TaskRepository{Data: map[string]*model.Task{
				tt.tsk.ID: tt.tsk,
			}}

			svc := testNewService(
				&ServiceConfig{TaskRepo: repo},
			)

			err := svc.DeleteTask(context.TODO(), tt.input)
			if err != nil {
				t.Fatal(err)
			}

			_, ok := repo.Data[tt.tsk.ID]
			c.Assert(ok, qt.Equals, false)
		})
	}
}
