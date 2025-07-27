package service

import "go-web/internal/core/ports"

type authService struct {
	store ports.Store
	cache ports.Cache
}

func NewAuthService(store ports.Store, cache ports.Cache) ports.AuthService {
	return &authService{store, cache}
}

func (a *authService) Register() {

}

func (a *authService) Login() {

}
