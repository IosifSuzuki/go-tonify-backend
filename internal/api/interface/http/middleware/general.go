package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
)

func abortWithResponse(ctx *gin.Context, statusCode int, err error) {
	errorMessage := err.Error()
	resp := dto.Response[any]{
		ErrorMessage: &errorMessage,
	}
	ctx.AbortWithStatusJSON(statusCode, resp)
}
