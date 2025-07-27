package hasher

import (
	"go-web/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewBcryptHasher() ports.Hasher {
	return &bcryptHasher{}
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *bcryptHasher) Compare(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
