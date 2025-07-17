package http

import (
	"encoding/json"
	"go-web/internal/core/models"
	"log/slog"
	"net/http"
)

type SuccessResponse[T any] struct {
	Data T         `json:"data"`
	Meta *MetaData `json:"meta,omitempty"`
}

type MetaData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type ErrorResponse struct {
	Error ErrorResponseDetail `json:"error"`
}

type ErrorResponseDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func respondSuccess[T any](w http.ResponseWriter, code int, data T, meta *MetaData) {
	writeJson(
		w,
		code,
		SuccessResponse[T]{
			Data: data,
			Meta: meta,
		})
}

func respondError(w http.ResponseWriter, err *models.AppError) {
	if err.IsInternal {
		slog.Error("internal error occurs", "error", err.InternalErr)
	}
	writeJson(
		w,
		err.StatusCode,
		ErrorResponse{
			Error: ErrorResponseDetail{
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
