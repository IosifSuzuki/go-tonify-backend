package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type MultipartForm struct {
	container container.Container
}

func NewMultipartForm(container container.Container) *MultipartForm {
	return &MultipartForm{
		container: container,
	}
}

func (m *MultipartForm) Limit(maxSize int64) gin.HandlerFunc {
	log := m.container.GetLogger()
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)
		if err := ctx.Request.ParseMultipartForm(maxSize); err != nil {
			log.Error("fail to parse multipart form", logger.FError(err))
			abortWithResponse(ctx, http.StatusRequestEntityTooLarge, err)
			return
		}
		ctx.Next()
	}
}
