package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	domain "go-web/internal/core/models"
	rest "go-web/internal/transport/http/models"
)

func respondSuccess(w http.ResponseWriter, code int, resp any) {
	writeJson(w, code, resp)
}

func respondError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.Internal(err)
	}
	if appErr.IsInternal {
		slog.Error("internal error occurs", "error", appErr.Err)
	}
	code := mapAppErrorTypeToStatusCode(appErr.Type)
	writeJson(
		w,
		code,
		&rest.ErrorResponseBody{
			StatusCode:   code,
			ErrorCode:    string(appErr.Type),
			ErrorMessage: appErr.Message,
		},
	)
}

func writeJson(w http.ResponseWriter, code int, target interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	//nolint:errcheck
	json.NewEncoder(w).Encode(target)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func mapAppErrorTypeToStatusCode(typ domain.ErrorType) int {
	switch typ {
	case domain.ErrInvalidParam, domain.ErrInvalidBody:
		return http.StatusBadRequest
	case domain.ErrInvalidAccess:
		return http.StatusUnauthorized
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrTooManyReq:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
