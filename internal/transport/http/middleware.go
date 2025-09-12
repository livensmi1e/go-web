package http

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-web/internal/core/models"
	"go-web/internal/platform"
)

func RegisterMiddlewares(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r
	for i := len(middlewares) - 1; i >= 0; i-- {
		s = middlewares[i](s)
	}
	return s
}

func LoggingMiddleware(next http.Handler) http.Handler {
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

func (h *apiHandler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.limiter == nil {
			next.ServeHTTP(w, r)
			return
		}
		ip := strings.Split(r.RemoteAddr, ":")[0]
		allowed, err := h.limiter.Allow(r.Context(), ip)
		if err != nil {
			respondError(w, models.Internal(err))
			return
		}
		if !allowed {
			respondError(w, models.TooManyRequests("Too many requests", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *apiHandler) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			respondError(w, models.InvalidAccess("Invalid or expired token", nil))
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := h.auth.Validate(tokenStr)
		if err != nil {
			respondError(w, models.InvalidAccess("Invalid or expired token", err))
			return
		}
		userId, ok := claims["sub"].(string)
		if !ok || userId == "" {
			respondError(w, models.InvalidAccess("Invalid token payload", nil))
			return
		}
		ctx := context.WithValue(r.Context(), platform.CtxUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HttpMetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		st := strconv.Itoa(sw.status)
		httpRequest.WithLabelValues(r.Method, st).Inc()
	})
}
