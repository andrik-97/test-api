package task

import (
	"os"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/payfazz/backend-template/pkg/txn/txnmock"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gopkg.in/go-playground/validator.v9"
)

// Creating a new service for testing purpose.
// It provides the default values for some fields.
func testNewService(cfg *ServiceConfig) Service {
	if cfg.Validate == nil {
		cfg.Validate = validator.New()
	}

	// In case cfg.RequestCount / cfg.RequestLatency is nil,
	// we'll use this registry to register them,
	// so the metrics don't collide with what's already in the
	// default registry.
	promRegistry := stdprometheus.NewRegistry()

	if cfg.RequestCount == nil {
		counterVec := stdprometheus.NewCounterVec(stdprometheus.CounterOpts{
			Namespace: "todo_api",
			Subsystem: "task_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"})
		promRegistry.Register(counterVec)
		cfg.RequestCount = kitprometheus.NewCounter(counterVec)
	}

	if cfg.RequestLatency == nil {
		summaryVec := stdprometheus.NewSummaryVec(stdprometheus.SummaryOpts{
			Namespace: "todo_api",
			Subsystem: "task_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"})
		promRegistry.Register(summaryVec)
		cfg.RequestLatency = kitprometheus.NewSummary(summaryVec)
	}

	if cfg.Logger == nil {
		cfg.Logger = log.NewLogfmtLogger(os.Stderr)
	}

	if cfg.Transactor == nil {
		cfg.Transactor = txnmock.NewTransactor()
	}

	return NewService(cfg)
}
