// Package http provides http handlers for task domain use cases.
package http

import (
	"github.com/payfazz/backend-template/internal/task"
)

// Server hosts endpoints to manage tasks.
type Server struct {
	taskSvc task.Service
}

// ServerConfig contains everything needed by a Server.
type ServerConfig struct {
	TaskSvc task.Service
}

// NewServer creates a new Server.
func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		taskSvc: cfg.TaskSvc,
	}
}
