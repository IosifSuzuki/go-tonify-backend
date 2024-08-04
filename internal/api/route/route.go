package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/service"
)

func Setup(gin *gin.Engine, authService service.AuthService, profileService service.ProfileService) {
	NewAuthRouter(gin.Group("auth"), authService)
	NewProfileRouter(gin.Group("profile"), profileService)
}
