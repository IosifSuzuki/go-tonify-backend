package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
)

func NewAuthRouter(group *gin.RouterGroup) {
	ac := &controller.AuthController{}
	group.POST("/sign-in", ac.SignIn)
	group.POST("/sign-up", ac.SignUp)
}
