package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
)

func Setup(
	gin *gin.Engine,
	container container.Container,
	authService service.AuthService,
	countryRepository repository.CountryRepository,
) {
	NewAuthRouter(gin.Group("auth"), authService)
	NewCommonRouter(gin.Group("common"), container, countryRepository)
	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
