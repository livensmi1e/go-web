package service

import "go-web/internal/core/ports"

type authService struct {
	store  ports.Store
	hasher ports.Hasher
}

func NewAuthService(store ports.Store, hasher ports.Hasher) ports.AuthService {
	return &authService{store, hasher}
}

func (a *authService) Register() {

}

func (a *authService) Login() {

}
