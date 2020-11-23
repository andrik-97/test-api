package model

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCompleteTasks(t *testing.T) {
	tests := map[string]struct {
		task *Task
		want *Task
	}{
		"given completed=false, success to complete the task": {
			task: &Task{Completed: false},
			want: &Task{Completed: true},
		},
		"given completed=true, success to complete the task": {
			task: &Task{Completed: true},
			want: &Task{Completed: true},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			tc.task.Complete()
			c.Assert(tc.task, qt.CmpEquals(cmpopts.IgnoreFields(Task{}, "CompletedAt")), tc.want)
			c.Assert(tc.task.CompletedAt, qt.Not(qt.IsNil))
		})
	}
}
