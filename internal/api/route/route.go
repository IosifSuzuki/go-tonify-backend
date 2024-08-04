package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-tonify-backend/internal/service"
)

func Setup(gin *gin.Engine, authService service.AuthService) {
	NewAuthRouter(gin.Group("auth"), authService)
	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
