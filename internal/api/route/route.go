package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/usecase"
)

func Setup(gin *gin.Engine, profileUseCase usecase.ProfileUseCase) {
	NewAuthRouter(gin.Group("auth"))
	NewProfileRouter(gin.Group("profile"), profileUseCase)
}
