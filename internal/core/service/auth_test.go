package service_test

import (
	"context"
	"go-web/internal/core/models"
	"go-web/internal/core/service"
	"go-web/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	t.Run("should register a new user", func(t *testing.T) {
		store := new(mocks.MockStore)
		hasher := new(mocks.MockHasher)
		token := new(mocks.MockToken)
		email := "user@test.com"
		password := "password"
		hashedPassword := "hashedPassword"
		store.On("FindByEmail", ctx, email).Return((*models.User)(nil), nil)
		hasher.On("Hash", password).Return(hashedPassword, nil)
		store.On("Create", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Email == email && u.PasswordHash == hashedPassword
		})).Return(&models.User{
			Email:        email,
			PasswordHash: hashedPassword,
		}, nil)
		authService := service.NewAuthService(store, hasher, token)
		user, err := authService.Register(ctx, email, password)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, hashedPassword, user.PasswordHash)
		store.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("should not register a user with existing email", func(t *testing.T) {
		store := new(mocks.MockStore)
		hasher := new(mocks.MockHasher)
		token := new(mocks.MockToken)
		email := "user@test.com"
		store.On("FindByEmail", ctx, email).Return(&models.User{
			Id:           "1",
			Email:        email,
			PasswordHash: "hashedPassword",
		}, nil)
		authService := service.NewAuthService(store, hasher, token)
		user, err := authService.Register(ctx, email, "password")
		assert.Error(t, err)
		assert.Nil(t, user)
		store.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	t.Run("should login a user and return a token", func(t *testing.T) {
		store := new(mocks.MockStore)
		hasher := new(mocks.MockHasher)
		token := new(mocks.MockToken)
		email := "user@test.com"
		password := "password"
		hashedPassword := "hashedPassword"
		expectedToken := "token"
		store.On("FindByEmail", ctx, email).Return(&models.User{
			Id:           "1",
			Email:        email,
			PasswordHash: hashedPassword,
		}, nil)
		hasher.On("Compare", hashedPassword, password).Return(nil)
		token.On("Generate", map[string]interface{}{
			"sub":   "1",
			"email": email,
		}).Return(expectedToken, nil)
		authService := service.NewAuthService(store, hasher, token)
		tok, err := authService.Login(ctx, email, password)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, tok)
		store.AssertExpectations(t)
		hasher.AssertExpectations(t)
		token.AssertExpectations(t)
	})
	t.Run("should not login with incorrect email", func(t *testing.T) {
		store := new(mocks.MockStore)
		hasher := new(mocks.MockHasher)
		token := new(mocks.MockToken)
		email := "wrong@test.com"
		password := "password"
		store.On("FindByEmail", ctx, email).Return((*models.User)(nil), nil)
		authService := service.NewAuthService(store, hasher, token)
		tok, err := authService.Login(ctx, email, password)
		assert.Error(t, err)
		assert.Empty(t, tok)
		store.AssertExpectations(t)
	})
	t.Run("should not login with incorrect password", func(t *testing.T) {
		store := new(mocks.MockStore)
		hasher := new(mocks.MockHasher)
		token := new(mocks.MockToken)

		email := "user@test.com"
		password := "wrongPassword"
		hashedPassword := "hashedPassword"
		store.On("FindByEmail", ctx, email).Return(&models.User{
			Id:           "1",
			Email:        email,
			PasswordHash: hashedPassword,
		}, nil)
		hasher.On("Compare", hashedPassword, password).Return(assert.AnError)
		authService := service.NewAuthService(store, hasher, token)
		tok, err := authService.Login(ctx, email, password)
		assert.Error(t, err)
		assert.Empty(t, tok)
		store.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})
}
