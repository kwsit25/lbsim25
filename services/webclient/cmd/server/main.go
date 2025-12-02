package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"schwarzit.load/services/webclient/internal/config"
	"schwarzit.load/services/webclient/internal/httpserver"
	"schwarzit.load/services/webclient/internal/metrics"
	"schwarzit.load/services/webclient/internal/service"
)

var logger *slog.Logger

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(); err != nil {
		logger.ErrorContext(context.Background(), "an error occurred", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, cfgErr := config.Get()
	if cfgErr != nil {
		os.Exit(1)
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	hostname, _ := os.Hostname()
	logger.With(slog.String("pod_name", hostname))

	logger.Info("starting load server")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	prometheusMetrics := metrics.NewMetrics()

	httpServer, err := httpserver.NewServer(
		fmt.Sprintf(":%d", cfg.HTTPPort),
		prometheusMetrics,
		logger,
	)
	if err != nil {
		return fmt.Errorf("creating http server: %w", err)
	}

	registry, err := prometheusMetrics.Registry()
	if err != nil {
		return fmt.Errorf("getting prometheus registry: %w", err)
	}

	httpServer.RegisterHTTP(
		http.MethodGet,
		"/metrics",
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			Registry: registry,
		}))

	// start here http server in secondary goroutine only for metrics and live & ready endpoints
	go runHTTPServer(ctx, logger, cfg, httpServer)

	srv := service.NewUserService(cfg, logger, prometheusMetrics)
	// start the service in another goroutine
	go srv.Start(ctx)

	return handleGracefulShutdown(ctx, logger)
}

func runHTTPServer(ctx context.Context, logger *slog.Logger, cfg *config.Config, httpServer *httpserver.Server) {
	logger.InfoContext(ctx, "listening on http", slog.Int("port", cfg.HTTPPort))

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.ErrorContext(ctx, "http server listen error", slog.String("error", err.Error()))
	}
}

func handleGracefulShutdown(
	ctx context.Context,
	logger *slog.Logger,
) error {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	<-signalChannel

	shutdownTimeout := 1 * time.Second

	logger.InfoContext(ctx, "shutting down", slog.Duration("timeout", shutdownTimeout))
	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	return nil
}
