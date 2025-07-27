package http

import (
	"encoding/json"
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
			Data: data,
			Meta: meta,
		})
}

func respondError(w http.ResponseWriter, err *domain.AppError) {
	if err.IsInternal {
		slog.Error("internal error occurs", "error", err.InternalErr)
	}
	writeJson(
		w,
		err.StatusCode,
		rest.ErrorResponse{
			Error: rest.ErrorResponseDetail{
				Type:    err.Type,
				Message: err.Message,
			},
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
