package models

type ErrorResponseBody struct {
	StatusCode   int    `json:"statusCode"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
