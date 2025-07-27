package store

import (
	"context"
	"go-web/internal/core/models"
)

func (p *pgStore) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, nil
}

func (p *pgStore) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
