package service

import "go-web/internal/core/ports"

type authService struct {
	store  ports.Store
	hasher ports.Hasher
	token  ports.TokenGenerator
}

func NewAuthService(store ports.Store, hasher ports.Hasher, token ports.TokenGenerator) ports.AuthService {
	return &authService{store, hasher, token}
}

func (a *authService) Register() {

}

func (a *authService) Login() {

}
