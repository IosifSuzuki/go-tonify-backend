package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/service"
)

func NewAuthRouter(group *gin.RouterGroup, authService service.AuthService) {
	ac := &controller.AuthController{AuthService: authService}

	group.POST("/sign-up", ac.AccountSignUp)
	group.POST("/sign-in", ac.AccountSignIn)
}
