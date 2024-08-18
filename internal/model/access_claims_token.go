package model

import "github.com/golang-jwt/jwt/v5"

type AccessClaimsToken struct {
	jwt.RegisteredClaims
	ID int64
}
