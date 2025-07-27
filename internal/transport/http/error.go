package http

import (
	"go-web/internal/core/models"
	"net/http"
)

func InvalidParam(msg string) *models.AppError {
	return &models.AppError{Type: models.ERR_INVALID_PARAMETER, Message: msg, StatusCode: http.StatusBadRequest}
}

func InvalidAccess(msg string) *models.AppError {
	return &models.AppError{Type: models.ERR_INVALID_ACCESS, Message: msg, StatusCode: http.StatusForbidden}
}

func UnknownError(inernalErr error) *models.AppError {
	return &models.AppError{Type: models.ERR_UNKNOWN, Message: models.MSG_UNKNOWN, StatusCode: http.StatusInternalServerError, IsInternal: true, InternalErr: inernalErr}
}
