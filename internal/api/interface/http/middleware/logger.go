package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/pkg/curl"
	"go-tonify-backend/pkg/logger"
)

type Logger struct {
	container container.Container
}

func NewLogger(container container.Container) *Logger {
	return &Logger{
		container: container,
	}
}

func (l *Logger) Logging() gin.HandlerFunc {
	log := l.container.GetLogger()
	return func(ctx *gin.Context) {
		curlCmd, err := curl.GetCurlCommand(ctx.Request)
		if err != nil {
			log.Error("fail to convert request to curl command", logger.FError(err))
			ctx.Next()
			return
		}
		log.Debug(curlCmd.String())
		ctx.Next()
	}
}
