package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/test-api/internal/task"
	taskhttp "github.com/payfazz/test-api/internal/task/http"
	taskpg "github.com/payfazz/test-api/internal/task/model/postgres"
	taskview "github.com/payfazz/test-api/internal/task/view"
	"github.com/payfazz/test-api/pkg/txn/txnsql"
	"github.com/payfazz/pkg/env"
	"github.com/pressly/goose"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"gopkg.in/go-playground/validator.v9"
)

var (
	port        = env.String("PORT", "8080")
	postgresURL = env.String("POSTGRES_URL", "postgres://postgres@localhost/todosvc?sslmode=disable")
	gooseDir    = env.String("GOOSE_DIR", "./internal/postgres/migration")
)

func main() {
	run(nil)
}

func run(
	// For tests, the technique is from: https://github.com/hashicorp/vault/blob/master/command/agent.go.
	startedCh chan (struct{}),
) {
	// 1. Create a single logger, which we'll use and give to other components.
	//

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 2. Migrate the db.
	//

	db, err := sql.Open("postgres", *postgresURL)
	if err != nil {
		level.Error(logger).Log("msg", "failed to open db", "err", err)
		os.Exit(1)
	}
	if err := goose.Up(db, *gooseDir); err != nil {
		level.Error(logger).Log("msg", "failed to goose up", "err", err)
		os.Exit(1)
	}
	if err := db.Close(); err != nil {
		level.Error(logger).Log("msg", "failed to close db", "err", err)
		os.Exit(1)
	}

	// 3. Create a db connection, which we'll pass to other components.
	//

	dbx := sqlx.MustOpen("postgres", *postgresURL)

	// 4. Initialize the http routes.
	//

	validate := validator.New()
	taskSrv := taskhttp.NewServer(&taskhttp.ServerConfig{
		TaskSvc: task.NewService(&task.ServiceConfig{
			TaskRepo: taskpg.NewTaskRepository(dbx),
			TaskView: taskview.NewTaskView(dbx, validate),
			Validate: validate,
			RequestCount: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "todo_api",
				Subsystem: "task_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, []string{"method"}),
			RequestLatency: kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "todo_api",
				Subsystem: "task_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, []string{"method"}),
			Transactor: txnsql.NewTransactor(dbx),
			Logger:     logger,
		}),
	})

	r := chi.NewRouter()

	// Basic CORS
	//
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "POST-Platform"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	// Mount application routes.
	r.Mount("/task", taskSrv.Routes())

	// Set the route for the prometheus metrics.
	//
	// In order for our infra prometheus server to scrape the metrics,
	// the app should expose :9100/metrics for the prometheus metrics.
	//
	go func() {
		r := chi.NewRouter()
		r.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9100", r)
	}()

	// 5. Run the app http server with gracefully shutting down capability.
	//
	// https://golang.org/pkg/net/http/#Server.Shutdown.
	//

	srv := http.Server{Addr: fmt.Sprintf(":%s", *port), Handler: r}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			level.Error(logger).Log("msg", "failed HTTP server shutdown", "err", err)
		}
		close(idleConnsClosed)
	}()

	logger.Log("msg", "server is running", "port", *port)

	// Inform any tests that the server is ready
	select {
	case startedCh <- struct{}{}:
	default:
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		level.Error(logger).Log("msg", "failed HTTP server ListenAndServe", "err", err)
		os.Exit(1)
	}

	<-idleConnsClosed
}
