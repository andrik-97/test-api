package http

import (
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/payfazz/go-httpjson/errencoder"
)

// Routes defines and returns the routes for this server.
func (s *Server) Routes() *chi.Mux {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errencoder.New()),
	}

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Post("/tasks", s.handleCreateTask(options...))
		r.Post("/tasks/{task_id}/complete", s.handleCompleteTask(options...))
		r.Delete("/tasks/{task_id}", s.handleDeleteTask(options...))
		r.Get("/tasks/{task_id}", s.handleViewTask(options...))
		r.Get("/tasks", s.handleViewTasks(options...))
	})

	return r
}
