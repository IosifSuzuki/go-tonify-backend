package route

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/controller"
	"go-tonify-backend/internal/api/middleware"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/repository"
)

func NewAccountRouter(
	group *gin.RouterGroup,
	container container.Container,
	accountRepository repository.AccountRepository,
	authMiddleware *middleware.Auth,
) {
	group.Use(authMiddleware.Authorization())
	account := controller.NewAccount(container, accountRepository)
	group.GET("/my", account.GetMyAccount)
	group.PATCH("/edit", account.EditMyAccount)
}
