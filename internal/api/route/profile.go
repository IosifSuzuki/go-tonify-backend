package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/service"
)

func NewProfileRouter(group *gin.RouterGroup, profileService service.ProfileService) {
	pc := &controller.ProfileController{profileService}
	
	group.GET("/:id", pc.GetProfile)
}
