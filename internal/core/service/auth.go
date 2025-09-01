package service

import (
	"context"
	"time"

	"go-web/internal/core/models"
	"go-web/internal/core/ports"
	"go-web/internal/shared"

	"github.com/google/uuid"
)

type authService struct {
	store  ports.Store
	cache  ports.Cache
	hasher ports.Hasher
	token  ports.TokenGenerator
}

func NewAuthService(store ports.Store, cache ports.Cache, hasher ports.Hasher, token ports.TokenGenerator) ports.AuthService {
	return &authService{store, cache, hasher, token}
}

func (a *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
	userDb, err := a.store.FindByEmail(ctx, email)
	if err != nil {
		return nil, models.Internal(err)
	}
	if userDb != nil {
		return nil, models.Conflict("Email already in use", nil)
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

func (a *authService) Login(ctx context.Context, email, password string) (*models.AuthTokens, error) {
	user, err := a.store.FindByEmail(ctx, email)
	if user == nil {
		return nil, models.InvalidAccess("Email or password is incorrect", err)
	}
	if err != nil {
		return nil, models.Internal(err)
	}
	if err := a.hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, models.InvalidAccess("Email or password is incorrect", err)
	}
	claims := map[string]interface{}{
		"sub":   user.Id,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
		"jti":   shared.RandString(8),
	}
	accessToken, err := a.token.Generate(claims)
	if err != nil {
		return nil, models.Internal(err)
	}
	refreshToken := shared.RandString(16)
	refreshUser := models.RefreshUser{
		Id:    user.Id,
		Email: user.Email,
	}
	err = a.cache.SetWithTTL(refreshToken, refreshUser, 60*60*24*7) // 7 days
	if err != nil {
		return nil, models.Internal(err)
	}
	return &models.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (*models.AuthTokens, error) {
	var refreshUser models.RefreshUser
	err := a.cache.Get(refreshToken, &refreshUser)
	if err != nil {
		return nil, models.InvalidAccess("Invalid refresh token", err)
	}
	if refreshUser.Id == "" {
		return nil, models.InvalidAccess("Invalid refresh token", nil)
	}
	newClaims := map[string]any{
		"sub":   refreshUser.Id,
		"email": refreshUser.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
		"jti":   shared.RandString(8),
	}
	newAccessToken, err := a.token.Generate(newClaims)
	if err != nil {
		return nil, models.Internal(err)
	}
	newRefreshToken := shared.RandString(16)
	err = a.cache.SetWithTTL(newRefreshToken, refreshUser, 60*60*24*7) // 7 days
	if err != nil {
		return nil, models.Internal(err)
	}
	err = a.cache.Delete(refreshToken)
	if err != nil {
		return nil, models.Internal(err)
	}
	return &models.AuthTokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	return a.cache.Delete(refreshToken)
}

func (a *authService) Validate(token string) (map[string]interface{}, error) {
	return a.token.Validate(token)
}
