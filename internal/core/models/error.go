package models

type ErrorType string

const (
	ErrInvalidParam  ErrorType = "INVALID_PARAMETER"
	ErrInvalidBody   ErrorType = "INVALID_BODY"
	ErrInvalidAccess ErrorType = "INVALID_ACCESS"
	ErrConflict      ErrorType = "CONFLICT"
	ErrNotFound      ErrorType = "NOT_FOUND"
	ErrUnknown       ErrorType = "UNKNOWN"
)

const (
	MsgUnknown = "Please contact our support team for details"
)

type AppError struct {
	Type       ErrorType
	Message    string
	Err        error
	IsInternal bool
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func newAppError(t ErrorType, msg string, err error, isInternal bool) *AppError {
	return &AppError{
		Type:       t,
		Message:    msg,
		Err:        err,
		IsInternal: isInternal,
	}
}

func InvalidParam(msg string, err error) *AppError {
	return newAppError(ErrInvalidParam, msg, err, false)
}

func InvalidBody(msg string, err error) *AppError {
	return newAppError(ErrInvalidBody, msg, err, false)
}

func InvalidAccess(msg string, err error) *AppError {
	return newAppError(ErrInvalidAccess, msg, err, false)
}

func Conflict(msg string, err error) *AppError {
	return newAppError(ErrConflict, msg, err, false)
}

func NotFound(msg string, err error) *AppError {
	return newAppError(ErrNotFound, msg, err, false)
}

func Internal(err error) *AppError {
	return newAppError(ErrUnknown, MsgUnknown, err, true)
}
