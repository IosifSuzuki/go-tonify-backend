package encrypt

import (
	"github.com/golang-jwt/jwt/v5"
	"go-tonify-backend/internal/container"
)

type JWT struct {
	container container.Container
}

func NewJWT(container container.Container) *JWT {
	return &JWT{
		container: container,
	}
}

func (j *JWT) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.container.GetJWTSecretKey()))
}
