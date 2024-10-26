package encrypt

import (
	"github.com/golang-jwt/jwt/v5"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/model"
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

func (j *JWT) ParseAccessToken(tokenText string) (*model.AccessClaimsToken, error) {
	var accessClaimsToken model.AccessClaimsToken
	token, err := jwt.ParseWithClaims(tokenText, &accessClaimsToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.container.GetJWTSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, model.JWTNotValidError
	}
	return &accessClaimsToken, nil
}

func (j *JWT) ParseRefreshToken(tokenText string) (*model.RefreshClaimsToken, error) {
	var refreshClaimsToken model.RefreshClaimsToken
	token, err := jwt.ParseWithClaims(tokenText, &refreshClaimsToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.container.GetJWTSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, model.JWTNotValidError
	}
	return &refreshClaimsToken, nil
}
