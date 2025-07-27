package token

import (
	"errors"
	"fmt"
	"go-web/internal/core/ports"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtGenerator struct {
	secret []byte
	exp    time.Duration
}

func NewJwtGenerator(secret string, exp time.Duration) ports.TokenGenerator {
	return &jwtGenerator{secret: []byte(secret), exp: exp}
}

func (j *jwtGenerator) Generate(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	maps.Copy(jwtClaims, claims)
	jwtClaims["exp"] = time.Now().Add(j.exp).Unix()
	jwtClaims["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)
	return token.SignedString(j.secret)
}

func (j *jwtGenerator) Validate(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
