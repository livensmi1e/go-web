package http

import (
	"encoding/json"
	"errors"
	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"net/http"

	domain "go-web/internal/core/models"
	rest "go-web/internal/transport/http/models"

	httpSwagger "github.com/swaggo/http-swagger"
)

type apiHandler struct {
	auth      ports.AuthService
	validator ports.Validator
	cache     ports.Cache
}

func newApiHandler(opts ...func(h *apiHandler)) *apiHandler {
	h := &apiHandler{}
	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *apiHandler) registerRoutes(mux *http.ServeMux) {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /example", h.helloWorld)
	apiMux.HandleFunc("GET /error", h.giveError)
	apiMux.HandleFunc("POST /auth/register", h.register)
	apiMux.HandleFunc("POST /auth/login", h.login)
	apiMux.Handle("GET /me", h.authorize(http.HandlerFunc(h.me)))
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mux.Handle("/docs/", httpSwagger.WrapHandler)
}

// helloWorld godoc
//
//	@Summary		Example Hello API
//	@Description	Returns a greeting message from the service
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	rest.SuccessResponse[string]
//	@Router			/example [get]
func (h *apiHandler) helloWorld(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		service.HelloWord(),
		nil,
	)
}

// giveError godoc
//
//	@Summary		Example Error API
//	@Description	Simulates an internal error response
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/error [get]
func (h *apiHandler) giveError(w http.ResponseWriter, r *http.Request) {
	respondError(w, domain.Internal(errors.New("db connection failed")))
}

// register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account with the provided email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.RegisterRequest							true	"User's credentials"
//	@Success		201		{object}	models.SuccessResponse[models.RegisterResponse]	"User created successfully"
//	@Failure		500		{object}	models.ErrorResponse							"Internal server error"
//	@Router			/auth/register [post]
func (h *apiHandler) register(w http.ResponseWriter, r *http.Request) {
	var req rest.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, domain.InvalidBody("invalid request body", err))
		return
	}
	user, err := h.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	res := &rest.RegisterResponse{
		Id:    user.Id,
		Email: user.Email,
	}
	respondSuccess(
		w,
		http.StatusCreated,
		res,
		nil,
	)
}

// login godoc
//
//	@Summary		Authenticate a user
//	@Description	Logs in a user with their email and password, returning a JWT on success.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.LoginRequest								true	"User's credentials for login"
//	@Success		200		{object}	models.SuccessResponse[models.LoginResponse]	"Authentication successful"
//	@Failure		400		{object}	models.ErrorResponse							"Invalid request body"
//	@Failure		401		{object}	models.ErrorResponse							"Invalid credentials"
//	@Failure		500		{object}	models.ErrorResponse							"Internal server error"
//	@Router			/auth/login [post]
func (h *apiHandler) login(w http.ResponseWriter, r *http.Request) {
	var req rest.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, domain.InvalidBody("invalid request body", err))
		return
	}
	token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	respondSuccess(
		w,
		http.StatusOK,
		&rest.LoginResponse{Token: token, Type: "Bearer"},
		nil,
	)
}

func (h *apiHandler) me(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		"Me",
		nil,
	)
}
