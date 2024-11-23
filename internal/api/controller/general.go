package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/utils"
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
func prepareAttachment(fileName string) *domain.Attachment {
	documentExt := utils.ExtFromFileName(fileName)
	documentUUID := uuid.NewString()
	documentName := fmt.Sprintf("%s.%s", documentUUID, documentExt)
	return &domain.Attachment{
		FileName: documentName,
		Status:   string(model.PendingAttachmentStatus),
	}
}
