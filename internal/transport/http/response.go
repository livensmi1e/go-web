package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	domain "go-web/internal/core/models"
	rest "go-web/internal/transport/http/models"
)

func respondSuccess[T any](w http.ResponseWriter, code int, data T, meta *rest.MetaData) {
	writeJson(
		w,
		code,
		rest.SuccessResponse[T]{
			Data:   data,
			Meta:   meta,
			Status: "success",
		})
}

func respondError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.Internal(err)
	}
	if appErr.IsInternal {
		slog.Error("internal error occurs", "error", appErr.Err)
	}
	writeJson(
		w,
		mapAppErrorTypeToStatusCode(appErr.Type),
		rest.ErrorResponse{
			Error: rest.ErrorResponseDetail{
				Type:    string(appErr.Type),
				Message: appErr.Message,
			},
			Status: "error",
		})
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
	default:
		return http.StatusInternalServerError
	}
}
