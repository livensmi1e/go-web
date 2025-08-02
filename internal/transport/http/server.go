package http

import (
	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"go-web/internal/infra/cache"
	"go-web/internal/infra/hasher"
	"go-web/internal/infra/store"
	"go-web/internal/infra/token"
	"go-web/internal/infra/validator"
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
	mux := http.NewServeMux()
	api := newApiHandler(func(a *apiHandler) {
		var s ports.Store
		var c ports.Cache
		var h ports.Hasher
		var t ports.TokenGenerator

		s = store.NewPgStore(cfg.StoreAddr())
		if s != nil {
			slog.Info("connected to db server on", "addr", cfg.StoreAddr())
		}
		if cfg.CacheEnabled {
			c = cache.NewGobCache(cfg.CacheAddr())
			slog.Info("connected to cache server on", "addr", cfg.CacheAddr())
		}
		h = hasher.NewBcryptHasher()
		t = token.NewJwtGenerator(cfg.JwtSecret, time.Minute*5)

		a.auth = service.NewAuthService(s, h, t)
		a.validator = validator.NewValidator()
		a.cache = c
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
