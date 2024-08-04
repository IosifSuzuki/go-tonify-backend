package model

import "github.com/golang-jwt/jwt/v5"

type RefreshClaimsToken struct {
	jwt.RegisteredClaims
	ID int64
}
