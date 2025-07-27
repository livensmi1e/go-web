package ports

import "go-web/internal/core/models"

type TokenGenerator interface {
	Generate(user *models.User) (string, error)
	Parse(token string) (*models.User, error)
}
