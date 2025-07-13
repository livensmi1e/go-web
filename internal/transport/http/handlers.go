package http

import (
	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"net/http"
)

type ApiHandler struct {
	store ports.Store
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
