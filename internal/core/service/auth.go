package service

import (
	"context"

	"go-web/internal/core/models"
	"go-web/internal/core/ports"

	"github.com/google/uuid"
)

type authService struct {
	store  ports.Store
	hasher ports.Hasher
	token  ports.TokenGenerator
}

func NewAuthService(store ports.Store, hasher ports.Hasher, token ports.TokenGenerator) ports.AuthService {
	return &authService{store, hasher, token}
}

func (a *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
	if _, err := a.store.FindByEmail(ctx, email); err != nil {
		return nil, models.Internal(err)
	}
	hashedPassword, err := a.hasher.Hash(password)
	if err != nil {
		return nil, models.Internal(err)
	}
	user := &models.User{
		Id:           uuid.NewString(),
		Email:        email,
		PasswordHash: hashedPassword,
	}
	if _, err := a.store.Create(ctx, user); err != nil {
		return nil, models.Internal(err)
	}
	return user, nil
}

func (a *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.store.FindByEmail(ctx, email)
	if user == nil {
		return "", models.InvalidAccess("Email or password is incorrect", err)
	}
	if err != nil {
		return "", models.Internal(err)
	}
	if err := a.hasher.Compare(user.PasswordHash, password); err != nil {
		return "", models.InvalidAccess("Email or password is incorrect", err)
	}
	claims := map[string]interface{}{
		"sub":   user.Id,
		"email": user.Email,
	}
	return a.token.Generate(claims)
}

func (a *authService) Validate(token string) (map[string]interface{}, error) {
	return a.token.Validate(token)
}
