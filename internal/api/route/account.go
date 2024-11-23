package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/api/middleware"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
)

func NewAccountRouter(
	group *gin.RouterGroup,
	container container.Container,
	accountRepository repository.AccountRepository,
	attachmentRepository repository.AttachmentRepository,
	companyRepository repository.CompanyRepository,
	attachmentService service.AttachmentService,
	authMiddleware *middleware.Auth,
) {
	group.Use(authMiddleware.Authorization())
	account := controller.NewAccount(container, accountRepository, attachmentRepository, companyRepository, attachmentService)
	group.GET("/my", account.GetMyAccount)
	group.PATCH("/edit", account.EditMyAccount)
	group.GET("/matching", account.MatchAccounts)
}
