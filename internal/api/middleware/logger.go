package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/pkg/logger"
	"net/http/httputil"
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
		requestDump, err := httputil.DumpRequest(ctx.Request, true)
		if err != nil {
			log.Debug("fail to get dump request", logger.FError(err))
		} else {
			log.Debug("receive body", logger.F("requestDump", string(requestDump)))
		}
		ctx.Next()
	}
}
