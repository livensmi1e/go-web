package ports

import (
	"context"
	"go-web/internal/core/models"
)

type Store interface {
	UserStore
}

type UserStore interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
