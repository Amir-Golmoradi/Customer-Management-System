package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/repository"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/service"
	httptransport "github.com/Amir-Golmoradi/Customer-Management-System/internal/customer/transport/http"
	generated "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/config"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/health"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/httpserver"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/platform/metrics"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/shared/idempotency"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Server *http.Server
}

func New(_ context.Context, cfg *config.Config, log *slog.Logger, pool *pgxpool.Pool) *App {
	healthHandler := health.NewHandler(pool)
	metricsCollector := metrics.NewCollector(func() metrics.DBStats {
		stats := pool.Stat()
		return metrics.DBStats{
			AcquireCount:         int64(stats.AcquireCount()),
			AcquiredConns:        int32(stats.AcquiredConns()),
			IdleConns:            int32(stats.IdleConns()),
			TotalConns:           int32(stats.TotalConns()),
			ConstructingConns:    int32(stats.ConstructingConns()),
			MaxConns:             int32(stats.MaxConns()),
			AcquireDurationNanos: stats.AcquireDuration().Nanoseconds(),
		}
	})

	queries := generated.New(pool)
	customerRepo := repository.NewPostgresRepository(queries, metricsCollector)
	customerService := service.New(customerRepo)
	customerHandler := httptransport.NewHandler(customerService, idempotency.NewStore())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", healthHandler.Live)
	mux.HandleFunc("GET /health/ready", healthHandler.Ready)
	mux.Handle("GET /metrics", metricsCollector.Handler())

	mux.HandleFunc("POST /v1/customers", customerHandler.CreateCustomer)
	mux.HandleFunc("GET /v1/customers", customerHandler.ListCustomers)
	mux.HandleFunc("GET /v1/customers/{id}", customerHandler.GetCustomerByID)

	handler := httpserver.Chain(
		mux,
		httpserver.RequestID,
		httpserver.AccessLog(log),
		httpserver.Metrics(metricsCollector),
	)

	return &App{Server: httpserver.New(handler, cfg.HTTP)}
}
