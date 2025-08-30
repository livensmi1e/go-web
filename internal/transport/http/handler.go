package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-web/internal/core/ports"
	"go-web/internal/core/service"

	domain "go-web/internal/core/models"
	rest "go-web/internal/transport/http/models"

	httpSwagger "github.com/swaggo/http-swagger"
)

type apiHandler struct {
	auth      ports.AuthService
	validator ports.Validator
	cache     ports.Cache
	limiter   ports.RateLimiter
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
//	@Success		200	{object}	models.ExampleResponseBody
//	@Router			/example [get]
func (h *apiHandler) helloWorld(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		&rest.ExampleResponseBody{Data: service.HelloWord(), StatusCode: http.StatusOK},
	)
}

// giveError godoc
//
//	@Summary		Example Error API
//	@Description	Simulates an internal error response
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	models.ErrorResponseBody
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
//	@Param			payload	body		models.RegisterRequestBody	true	"User's credentials"
//	@Success		201		{object}	models.RegisterResponseBody	"User created successfully"
//	@Failure		500		{object}	models.ErrorResponseBody	"Internal server error"
//	@Router			/auth/register [post]
func (h *apiHandler) register(w http.ResponseWriter, r *http.Request) {
	var req rest.RegisterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, domain.InvalidBody("invalid request body", err))
		return
	}
	user, err := h.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	data := &rest.RegisterResponse{
		Id:    user.Id,
		Email: user.Email,
	}
	resp := &rest.RegisterResponseBody{
		Data:       data,
		StatusCode: http.StatusCreated,
	}
	respondSuccess(
		w,
		http.StatusCreated,
		resp,
	)
}

// login godoc
//
//	@Summary		Authenticate a user
//	@Description	Logs in a user with their email and password, returning a JWT on success.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.LoginRequestBody		true	"User's credentials for login"
//	@Success		200		{object}	models.LoginResponseBody	"Authentication successful"
//	@Failure		400		{object}	models.ErrorResponseBody	"Invalid request body"
//	@Failure		401		{object}	models.ErrorResponseBody	"Invalid credentials"
//	@Failure		500		{object}	models.ErrorResponseBody	"Internal server error"
//	@Router			/auth/login [post]
func (h *apiHandler) login(w http.ResponseWriter, r *http.Request) {
	var req rest.LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, domain.InvalidBody("invalid request body", err))
		return
	}
	token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	data := &rest.LoginResponse{Token: token, Type: "Bearer"}
	resp := &rest.LoginResponseBody{
		Data:       data,
		StatusCode: http.StatusOK,
	}
	respondSuccess(
		w,
		http.StatusOK,
		resp,
	)
}

func (h *apiHandler) me(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		&rest.GetMeResponseBody{Data: "me", StatusCode: http.StatusOK},
	)
}
