package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
)

func NewAuthRouter(
	group *gin.RouterGroup,
	container container.Container,
	authService service.AuthService,
	attachmentService service.AttachmentService,
	companyRepository repository.CompanyRepository,
) {
	ac := controller.NewAuthController(container, authService, attachmentService, companyRepository)

	group.POST("/sign-up", ac.AccountSignUp)
	group.POST("/sign-in", ac.AccountSignIn)
}
