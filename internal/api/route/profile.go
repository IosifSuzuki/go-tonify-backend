package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/usecase"
)

func NewProfileRouter(group *gin.RouterGroup, profileUseCase usecase.ProfileUseCase) {
	pc := &controller.ProfileController{profileUseCase}
	group.GET("/:id", pc.GetProfile)
}
