package controller

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/model"
	"net/http"
)

func sendError(ctx *gin.Context, err error, code int) {
	var response model.ErrorResponse
	response.Message = err.Error()
	ctx.JSON(code, response)
}

func sendResponseWithStatus(ctx *gin.Context, resp any, code int) {
	ctx.JSON(code, resp)
}

func sendResponse(ctx *gin.Context, resp any) {
	sendResponseWithStatus(ctx, resp, http.StatusOK)
}
