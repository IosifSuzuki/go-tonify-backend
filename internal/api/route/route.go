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
	authMiddleware *middleware.Auth,
	corsMiddleware *middleware.CORS,
	loggerMiddleware *middleware.Logger,
	accountRepository repository.AccountRepository,
	countryRepository repository.CountryRepository,
) {
	gin.Use(corsMiddleware.CORS())
	gin.Use(loggerMiddleware.Logging())
	NewAuthRouter(gin.Group("auth"), authService)
	NewCommonRouter(gin.Group("common"), container, countryRepository)
	NewAccountRouter(gin.Group("account"), container, accountRepository, authMiddleware)

	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
