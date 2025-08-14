package models

type ErrorResponseBody struct {
	StatusCode   int    `json:"status_code"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
