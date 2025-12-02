package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"schwarzit.load/services/webclient/internal/config"
	"schwarzit.load/services/webclient/internal/metrics"
)

type Service struct {
	logger     *slog.Logger
	metrics    *metrics.Metrics
	cfg        *config.Config
	httpClient *http.Client
}

func NewUserService(cfg *config.Config, logger *slog.Logger, metrics *metrics.Metrics) *Service {
	return &Service{
		logger:  logger,
		metrics: metrics,
		cfg:     cfg,
		httpClient: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
}

func (s *Service) Start(ctx context.Context) {
	s.logger.InfoContext(ctx, "starting web client service")

	if s.cfg.RequestCount <= 0 || s.cfg.RequestInterval <= 0 {
		s.logger.ErrorContext(ctx, "invalid request config",
			slog.Int("request_count", s.cfg.RequestCount),
			slog.Int("request_interval", s.cfg.RequestInterval),
		)
		return
	}

	tickDuration := time.Duration(s.cfg.RequestInterval) * time.Millisecond
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.InfoContext(ctx, "stopping web client service")
			return
		case <-ticker.C:
			for i := 0; i < s.cfg.RequestCount; i++ {
				go s.performRequest(ctx)
			}
		}
	}
}

func (s *Service) performRequest(ctx context.Context) {
	hostname, _ := os.Hostname()
	targetUrl := fmt.Sprintf("%s?source=%s&mode=%s", s.cfg.TargetURL, hostname, s.cfg.Mode)
	s.logger.InfoContext(ctx, "call to url", slog.String("url", targetUrl))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to create request", slog.String("error", err.Error()))
		return
	}

	s.metrics.IncHttpSend(hostname)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to perform request", slog.String("error", err.Error()))
		// s.metrics.IncHttpRequestError("webclient") // TODO: wanted?
		return
	}
	defer resp.Body.Close()

	s.logger.InfoContext(ctx, "request successful", slog.String("status", resp.Status))
	_, _ = io.Copy(io.Discard, resp.Body)
}
