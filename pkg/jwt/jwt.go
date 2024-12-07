package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"go-tonify-backend/pkg/jwt/model"
)

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		secretKey: secretKey,
	}
}

func (j *JWT) GeneratePairJWT(accountID int64) (*model.PairToken, error) {
	accessClaimsToken := model.AccessClaimsToken{
		ID: accountID,
	}
	refreshClaimsToken := model.RefreshClaimsToken{
		ID: accountID,
	}
	accessToken, err := j.generateToken(accessClaimsToken)
	if err != nil {
		return nil, model.CreationJWTError
	}
	refreshToken, err := j.generateToken(refreshClaimsToken)
	if err != nil {
		return nil, model.CreationJWTError
	}
	return &model.PairToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (j *JWT) ParseAccessToken(tokenText string) (*model.AccessClaimsToken, error) {
	var accessClaimsToken model.AccessClaimsToken
	token, err := jwt.ParseWithClaims(tokenText, &accessClaimsToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
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
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, model.JWTNotValidError
	}
	return &refreshClaimsToken, nil
}

func (j *JWT) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
