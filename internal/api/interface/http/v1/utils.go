package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/utils"
)

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

func successResponse[T any](ctx *gin.Context, code int, model T) {
	var response = dto.Response{
		Response: &model,
	}
	ctx.JSON(code, response)
}

func failResponse(ctx *gin.Context, code int, error error, internalError error) {
	var errorMessage = error.Error()
	var response = dto.Response{
		Response:     nil,
		ErrorMessage: &errorMessage,
	}
	if internalError != nil {
		response.ErrorCause = utils.NewString(internalError.Error())
	}
	ctx.JSON(code, response)
}
