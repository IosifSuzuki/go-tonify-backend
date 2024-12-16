package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
)

func abortWithResponse(ctx *gin.Context, statusCode int, err error) {
	errorMessage := err.Error()
	resp := dto.Response{
		ErrorMessage: &errorMessage,
	}
	ctx.AbortWithStatusJSON(statusCode, resp)
}

func getAccountID(ctx *gin.Context) (*int64, error) {
	accountIDValue, exist := ctx.Get(dto.AccountIDKey)
	if !exist {
		return nil, dto.MissingAccountIDError
	}
	if accountIDValue == nil {
		return nil, dto.NilError
	}
	accountID, ok := accountIDValue.(int64)
	if !ok {
		return nil, dto.CastTypeError
	}
	return &accountID, nil
}
