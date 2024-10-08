package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/model"
)

func sendResponse(ctx *gin.Context, statusCode int, err error) {
	resp := model.ErrorResponse{
		Message: err.Error(),
	}
	ctx.AbortWithStatusJSON(statusCode, resp)
}
