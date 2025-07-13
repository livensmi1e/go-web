package http

import (
	"log/slog"
	"net/http"
	"time"
)

func registerMiddlewares(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r
	for i := len(middlewares) - 1; i >= 0; i-- {
		s = middlewares[i](s)
	}
	return s
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"duration_ms", duration,
		)
	})
}
