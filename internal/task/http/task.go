package http

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/payfazz/test-api/internal/task"
	"github.com/payfazz/go-httpjson"
)

func (s *Server) handleViewTask(options ...kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			input := request.(*task.ViewTaskInput)
			res, err := s.taskSvc.ViewTask(ctx, input)
			return res, err
		},
		httpjson.MakeDecode(&task.ViewTaskInput{}),
		httpjson.MakeEncode(),
		options...,
	).ServeHTTP
}

func (s *Server) handleViewTasks(options ...kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			input := request.(*task.ViewTasksInput)
			res, err := s.taskSvc.ViewTasks(ctx, input)
			return res, err
		},
		httpjson.MakeDecode(&task.ViewTasksInput{}),
		httpjson.MakeEncode(),
		options...,
	).ServeHTTP
}

func (s *Server) handleCreateTask(options ...kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			input := request.(*task.CreateTaskInput)
			res, err := s.taskSvc.CreateTask(ctx, input)
			return res, err
		},
		httpjson.MakeDecode(&task.CreateTaskInput{}),
		httpjson.MakeEncode(),
		options...,
	).ServeHTTP
}

func (s *Server) handleCompleteTask(options ...kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			input := request.(*task.CompleteTaskInput)
			err := s.taskSvc.CompleteTask(ctx, input)
			return nil, err
		},
		httpjson.MakeDecode(&task.CompleteTaskInput{}),
		httpjson.MakeEncode(),
		options...,
	).ServeHTTP
}

func (s *Server) handleDeleteTask(options ...kithttp.ServerOption) http.HandlerFunc {
	return kithttp.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			input := request.(*task.DeleteTaskInput)
			err := s.taskSvc.DeleteTask(ctx, input)
			return nil, err
		},
		httpjson.MakeDecode(&task.DeleteTaskInput{}),
		httpjson.MakeEncode(),
		options...,
	).ServeHTTP
}
