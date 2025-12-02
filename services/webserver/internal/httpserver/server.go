package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.opentelemetry.io/otel/trace"
	"schwarzit.load/services/webserver/internal/config"
	"schwarzit.load/services/webserver/internal/metrics"
)

type Server struct {
	server  http.Server
	mux     *http.ServeMux
	metrics *metrics.Metrics
	logger  *slog.Logger
}

const (
	defaultReadTimeout  = 2500 * time.Millisecond
	defaultWriteTimeout = 5000 * time.Millisecond
)

func NewServer(addr string, m *metrics.Metrics, logger *slog.Logger, cfg *config.Config) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ready", ready)
	mux.HandleFunc("GET /live", live)

	loadApi := setupLoadApi(logger, m, cfg)
	mux.HandleFunc("GET /api/load", loadApi)

	tracingMiddleware := otelhttp.NewMiddleware("http",
		otelhttp.WithFilter(func(request *http.Request) bool {
			return request.URL.Path != "/ready" && request.URL.Path != "/live" && request.URL.Path != "/metrics"
		}),
	)

	return &Server{
		server: http.Server{
			Handler:      tracingMiddleware(urlAttributesMiddleware(mux)),
			Addr:         addr,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		mux:     mux,
		metrics: m,
		logger:  logger,
	}, nil
}

func (s *Server) RegisterHTTP(method, path string, handler http.Handler) {
	s.mux.Handle(method+" "+path, handler)
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func ready(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func setupLoadApi(logger *slog.Logger, metrics *metrics.Metrics, cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		// get query parameter
		sourceQuery := req.URL.Query().Get("source")
		if sourceQuery == "" {
			sourceQuery = "unknown"
		}
		modeQuery := req.URL.Query().Get("mode")
		if modeQuery == "" {
			modeQuery = "unknown"
		}
		logger.Info("Handling API request",
			slog.String("url", req.URL.String()),
			slog.String("source", sourceQuery),
			slog.String("mode", modeQuery),
		)
		metrics.IncHttpRequestCount(sourceQuery, modeQuery)
		time.Sleep(time.Duration(cfg.ResponseDelay) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}
}

func urlAttributesMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		span := trace.SpanFromContext(request.Context())
		span.SetAttributes(semconv.URLPath(request.URL.Path))

		if request.URL.RawQuery != "" {
			span.SetAttributes(semconv.URLQuery(request.URL.RawQuery))
		}

		h.ServeHTTP(writer, request)
	})
}
