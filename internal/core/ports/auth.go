package ports

import (
	"context"
	"go-web/internal/core/models"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	Validate(token string) (map[string]interface{}, error)
}
