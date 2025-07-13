package http

import (
	"go-web/internal/infra/cache"
	"go-web/internal/infra/store"
	"go-web/internal/platform"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func newServer(opts ...func(*http.Server)) *http.Server {
	s := &http.Server{}
	for _, o := range opts {
		o(s)
	}
	return s
}

func withAddr(addr string) func(*http.Server) {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

func withTimeouts(read, write, header time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadTimeout = read
		s.WriteTimeout = write
		s.ReadHeaderTimeout = header
	}
}

func withHandler(h http.Handler) func(*http.Server) {
	return func(s *http.Server) {
		s.Handler = h
	}
}

func RunServer(cfg *platform.Config) error {
	// TODO: refine log placement
	mux := http.NewServeMux()
	api := newApiHandler(func(h *ApiHandler) {
		h.store = store.NewPg(cfg.StoreAddr())
		if cfg.CacheEnable {
			slog.Info("connected to cache server on", "addr", cfg.CacheAddr())
			h.cache = cache.NewGobCache(cfg.CacheAddr())
		}
	})
	api.registerRoutes(mux)
	handler := registerMiddlewares(mux, loggingMiddleware)
	server := newServer(
		withAddr(cfg.HttpServerAddr()),
		withHandler(cors.Default().Handler(handler)),
		withTimeouts(5*time.Second, 10*time.Second, 2*time.Second),
	)
	return server.ListenAndServe()
}
