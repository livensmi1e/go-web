package http

import (
	"encoding/json"
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
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	//nolint:errcheck
	json.NewEncoder(w).Encode(SuccessResponse[T]{
		Data: data,
		Meta: meta,
	})
}

// func respondError(w http.ResponseWriter, code int, err models.Error) {
// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(ErrorResponse{
// 		Error: ErrorResponseDetail{
// 			Type:    err.HtppErr().Error(),
// 			Message: err.AppErr().Error(),
// 		},
// 	})
// }
