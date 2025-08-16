package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go-web/internal/core/models"
	"go-web/internal/platform"
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
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sw, r)
		duration := time.Since(start)
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"remote", r.RemoteAddr,
			"duration_ms", duration,
		)
	})
}

func (h *apiHandler) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			respondError(w, models.InvalidAccess("invalid or expired token", nil))
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := h.auth.Validate(tokenStr)
		if err != nil {
			respondError(w, models.InvalidAccess("invalid or expired token", err))
			return
		}
		userId, ok := claims["sub"].(string)
		if !ok || userId == "" {
			respondError(w, models.InvalidAccess("invalid token payload", nil))
			return
		}
		ctx := context.WithValue(r.Context(), platform.CtxUserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
