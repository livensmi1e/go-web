package http

import (
	"errors"
	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"net/http"
)

type ApiHandler struct {
	store ports.Store
	cache ports.Cache
}

func newApiHandler(opts ...func(h *ApiHandler)) *ApiHandler {
	h := &ApiHandler{}
	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *ApiHandler) registerRoutes(mux *http.ServeMux) {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /example", h.HelloWorld)
	apiMux.HandleFunc("GET /error", h.GiveError)
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}

func (h *ApiHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		service.HelloWord(),
		nil,
	)
}

func (h *ApiHandler) GiveError(w http.ResponseWriter, r *http.Request) {
	respondError(w, UnknownError(errors.New("db connection failed")))
}
