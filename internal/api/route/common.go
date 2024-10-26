package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
)

func NewCommonRouter(group *gin.RouterGroup, container container.Container, countryRepository repository.CountryRepository) {
	c := &controller.CommonController{
		Container:         container,
		CountryRepository: countryRepository,
	}
	group.GET("/countries", c.Countries)
}
