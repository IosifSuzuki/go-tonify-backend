package v1

import (
	"github.com/gin-gonic/gin"
	v "github.com/go-playground/validator/v10"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/utils"
	"net/http"
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

func failResponse(ctx *gin.Context, code int, err error, internalError error) {
	var internalErrorMsg *string
	if internalError != nil {
		internalErrorMsg = utils.NewString(internalError.Error())
	}
	sendFailResponse(ctx, code, err.Error(), internalErrorMsg)
}

func badRequestResponse(ctx *gin.Context, validator validator.HttpValidator, err error, internalError error) {
	if validationErrors, ok := internalError.(v.ValidationErrors); ok {
		internalErrorMsg := validator.Translate(validationErrors)
		sendFailResponse(ctx, http.StatusBadRequest, err.Error(), &internalErrorMsg)
		return
	}
	failResponse(ctx, http.StatusBadRequest, err, internalError)
}

func sendFailResponse(ctx *gin.Context, code int, errMsg string, internalErrorMsg *string) {
	var response = dto.Response{
		Response:     nil,
		ErrorMessage: &errMsg,
	}
	if internalErrorMsg != nil {
		response.ErrorCause = internalErrorMsg
	}
	ctx.JSON(code, response)
}
