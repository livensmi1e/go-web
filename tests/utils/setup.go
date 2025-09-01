package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-web/internal/core/service"
	"go-web/internal/infra/cache"
	"go-web/internal/infra/hasher"
	"go-web/internal/infra/limiter"
	"go-web/internal/infra/store"
	"go-web/internal/infra/token"
	"go-web/internal/infra/validator"
	httpTransport "go-web/internal/transport/http"

	"github.com/stretchr/testify/require"
)

type TestServer struct {
	Server *httptest.Server
	Client *http.Client
}

func SetupTestServer() *TestServer {
	s := store.NewPgStore("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	c := cache.NewMemCache("localhost:11211")
	h := hasher.NewBcryptHasher()
	t := token.NewJwtGenerator("test_secret", time.Minute*5)
	l := limiter.NewMemLimiter(10, 30)
	v := validator.NewValidator()
	auth := service.NewAuthService(s, c, h, t)
	api := httpTransport.NewApiHandler(auth, v, c, l)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)
	handler := httpTransport.RegisterMiddlewares(mux, httpTransport.LoggingMiddleware, api.RateLimitMiddleware)
	ts := httptest.NewServer(handler)
	return &TestServer{
		Server: ts,
		Client: ts.Client(),
	}
}

func (ts *TestServer) DoRequest(t *testing.T, method, path string, body any, token string, respTarget any, wantStatus int, cookies ...*http.Cookie) *http.Response {
	var buf []byte
	var err error
	if body != nil {
		buf, err = json.Marshal(body)
		require.NoError(t, err)
	}
	req, err := http.NewRequest(method, ts.Server.URL+path, bytes.NewReader(buf))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	res, err := ts.Client.Do(req)
	require.NoError(t, err)
	//nolint:errcheck
	defer res.Body.Close()
	require.Equal(t, wantStatus, res.StatusCode)
	if respTarget != nil {
		err = json.NewDecoder(res.Body).Decode(respTarget)
		require.NoError(t, err)
	}
	return res
}

func GenUserEmail() string {
	return "user" + time.Now().Format("20060102150405") + "@test.com"
}
