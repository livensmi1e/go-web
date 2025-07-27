package models

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
