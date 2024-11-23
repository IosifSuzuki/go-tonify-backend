package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-tonify-backend/internal/api/middleware"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
)

func Setup(
	gin *gin.Engine,
	container container.Container,
	authService service.AuthService,
	attachmentService service.AttachmentService,
	authMiddleware *middleware.Auth,
	corsMiddleware *middleware.CORS,
	loggerMiddleware *middleware.Logger,
	accountRepository repository.AccountRepository,
	attachmentRepository repository.AttachmentRepository,
	companyRepository repository.CompanyRepository,
	countryRepository repository.CountryRepository,
) {
	gin.Use(corsMiddleware.CORS())
	gin.Use(loggerMiddleware.Logging())
	NewAuthRouter(gin.Group("auth"), container, authService, attachmentService, attachmentRepository, companyRepository)
	NewCommonRouter(gin.Group("common"), container, countryRepository)
	NewAccountRouter(gin.Group("account"), container, accountRepository, attachmentRepository, companyRepository, attachmentService, authMiddleware)

	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
