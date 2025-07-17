package models

const (
	ERR_INVALID_PARAMETER = "INVALID_PARAMETER"
	ERR_INVALID_ACCESS    = "INVALID_ACCESS"
	ERR_UNKNOWN           = "UNKNOWN_ERROR"
	MSG_UNKNOWN           = "Please contact our support team for details"
)

type AppError struct {
	Type        string
	Message     string
	StatusCode  int
	IsInternal  bool
	InternalErr error
}
