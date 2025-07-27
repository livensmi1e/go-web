package ports

import "go-web/internal/core/models"

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hashed string, plain string) bool
}

type TokenGenerator interface {
	Generate(user *models.User) (string, error)
	Parse(token string) (*models.User, error)
}
