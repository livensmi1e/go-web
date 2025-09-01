package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-web/internal/core/ports"
	"go-web/internal/core/service"
	"go-web/internal/shared"

	domain "go-web/internal/core/models"
	rest "go-web/internal/transport/http/models"

	httpSwagger "github.com/swaggo/http-swagger"
)

type apiHandler struct {
	auth      ports.AuthService
	validator ports.Validator
	cache     ports.Cache
	limiter   ports.RateLimiter

	env string
}

// This constructor is for test purpose only
func NewApiHandler(auth ports.AuthService, validator ports.Validator, cache ports.Cache, limiter ports.RateLimiter) *apiHandler {
	return &apiHandler{
		auth:      auth,
		validator: validator,
		cache:     cache,
		limiter:   limiter,
	}
}

func newApiHandler(opts ...func(h *apiHandler)) *apiHandler {
	h := &apiHandler{}
	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *apiHandler) RegisterRoutes(mux *http.ServeMux) {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /example", h.helloWorld)
	apiMux.HandleFunc("GET /error", h.giveError)
	apiMux.HandleFunc("POST /auth/register", h.register)
	apiMux.HandleFunc("POST /auth/login", h.login)
	apiMux.Handle("POST /auth/refresh", http.HandlerFunc(h.refresh))
	apiMux.Handle("POST /auth/logout", h.authorize(http.HandlerFunc(h.logout)))
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
		respondError(w, domain.InvalidBody("Invalid request body", err))
		return
	}
	if err := h.validator.Validate(req); err != nil {
		respondError(w, domain.InvalidBody("Invalid request body", err))
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
		respondError(w, domain.InvalidBody("Invalid request body", err))
		return
	}
	if err := h.validator.Validate(req); err != nil {
		respondError(w, domain.InvalidBody("Invalid request body", err))
	}
	tokens, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, err)
		return
	}
	data := &rest.LoginResponse{Token: tokens.AccessToken, Type: "Bearer"}
	resp := &rest.LoginResponseBody{
		Data:       data,
		StatusCode: http.StatusOK,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/auth/refresh",
		Secure:   !shared.IsDevelopmentEnv(h.env),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	respondSuccess(
		w,
		http.StatusOK,
		resp,
	)
}

// refresh godoc
//
//	@Summary		Refresh JWT tokens
//	@Description	Refreshes the JWT tokens using a valid refresh token
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	models.RefreshTokenResponseBody	"Token refreshed successfully"
//	@Failure		400	{object}	models.ErrorResponseBody			"Invalid request"
//	@Failure		401	{object}	models.ErrorResponseBody			"Invalid or expired refresh token"
//	@Failure		500	{object}	models.ErrorResponseBody			"Internal server error"
//	@Router			/auth/refresh [post]
func (h *apiHandler) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if refreshToken == nil || err != nil {
		respondError(w, domain.InvalidAccess("refresh token is required", nil))
		return
	}
	tokens, err := h.auth.Refresh(r.Context(), refreshToken.Value)
	if err != nil {
		respondError(w, err)
		return
	}
	data := &rest.LoginResponse{Token: tokens.AccessToken, Type: "Bearer"}
	resp := &rest.RefreshTokenResponseBody{
		Data:       data,
		StatusCode: http.StatusOK,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/auth/refresh",
		Secure:   !shared.IsDevelopmentEnv(h.env),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	respondSuccess(
		w,
		http.StatusOK,
		resp,
	)
}

// logout godoc
//
//	@Summary		Logout a user
//	@Description	Logs out a user by invalidating their refresh token
//	@Tags			Auth
//	@Produce		json
//	@Param			payload	body		models.LogoutRequestBody	true	"Refresh token to invalidate"
//	@Success		200		{object}	models.LogoutResponseBody	"Logout successful"
//	@Failure		400		{object}	models.ErrorResponseBody	"Invalid request body"
//	@Failure		500		{object}	models.ErrorResponseBody	"Internal server error"
//	@Router			/auth/logout [post]
func (h *apiHandler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if refreshToken == nil || err != nil {
		respondError(w, domain.InvalidAccess("refresh token is required", nil))
		return
	}
	if err := h.auth.Logout(r.Context(), refreshToken.Value); err != nil {
		respondError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   !shared.IsDevelopmentEnv(h.env),
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	})
	respondSuccess(
		w,
		http.StatusOK,
		&rest.LogoutResponseBody{Data: nil, StatusCode: http.StatusOK},
	)
}

func (h *apiHandler) me(w http.ResponseWriter, r *http.Request) {
	respondSuccess(
		w,
		http.StatusOK,
		&rest.GetMeResponseBody{Data: "me", StatusCode: http.StatusOK},
	)
}
