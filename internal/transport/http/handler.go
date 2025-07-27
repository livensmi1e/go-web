package http

import (
	"errors"
	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: Refactor this
type ApiHandler struct {
	auth      ports.AuthService
	validator ports.Validator
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
	mux.Handle("/docs/", httpSwagger.WrapHandler)
}

// HelloWorld godoc
//
//	@Summary		Example Hello API
//	@Description	Returns a greeting message from the service
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse[string]
//	@Router			/example [get]
func (h *ApiHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		service.HelloWord(),
		nil,
	)
}

// GiveError godoc
//
//	@Summary		Example Error API
//	@Description	Simulates an internal error response
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/error [get]
func (h *ApiHandler) GiveError(w http.ResponseWriter, r *http.Request) {
	respondError(w, UnknownError(errors.New("db connection failed")))
}
