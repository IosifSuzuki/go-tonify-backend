package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-tonify-backend/internal/api/interface/http/middleware"
	"go-tonify-backend/internal/api/interface/http/v1"
	"go-tonify-backend/internal/container"
	accountUsecase "go-tonify-backend/internal/domain/account/usecase"
	countryUsecase "go-tonify-backend/internal/domain/country/usecase"
)

type Handler struct {
	container      container.Container
	accountUsecase accountUsecase.Account
	countryUsecase countryUsecase.Country
}

func NewHandler(
	container container.Container,
	accountUsecase accountUsecase.Account,
	countryUsecase countryUsecase.Country,
) *Handler {
	return &Handler{
		container:      container,
		accountUsecase: accountUsecase,
		countryUsecase: countryUsecase,
	}
}

func (h *Handler) Run() error {
	r := gin.Default()

	loggerMiddleware := middleware.NewLogger(h.container)
	corsMiddleware := middleware.NewCORS(h.container)
	authMiddleware := middleware.NewAuth(h.container, h.accountUsecase)

	r.Use(corsMiddleware.CORS())
	r.Use(loggerMiddleware.Logging())

	v1 := r.Group("api/v1")

	authHandler := h.composeAuthHandler()

	authGroup := v1.Group("auth")
	{
		authGroup.POST("/sign-up", authHandler.SignUp)
		authGroup.POST("/sign-in", authHandler.SignIn)
	}
	accountHandler := h.composeAccount()
	accountGroup := v1.Group("account")
	accountGroup.Use(authMiddleware.Authorization())
	{
		accountGroup.GET("/my", accountHandler.GetMy)
		accountGroup.PATCH("/edit", accountHandler.EditMy)
		accountGroup.GET("/matching", accountHandler.MatchAccounts)
	}
	commonHandler := h.composeCommon()
	commonGroup := v1.Group("/common")
	{
		commonGroup.GET("/ping", commonHandler.Ping)
		commonGroup.GET("/countries", commonHandler.Countries)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go h.runHttpServer(r)
	return h.runHttpsServer(r)
}

func (h *Handler) runHttpServer(r *gin.Engine) error {
	return r.Run(
		h.container.GetServerConfig().Address(),
	)
}

func (h *Handler) runHttpsServer(r *gin.Engine) error {
	return r.RunTLS(
		h.container.GetServerConfig().SecureAddress(),
		"tls/public.pem",
		"tls/private.key",
	)
}

func (h *Handler) composeAuthHandler() *v1.AuthHandler {
	return v1.NewAuthHandler(h.container, h.accountUsecase)
}

func (h *Handler) composeAccount() *v1.AccountHandler {
	return v1.NewAccountHandler(h.container, h.accountUsecase)
}

func (h *Handler) composeCommon() *v1.CommonHandler {
	return v1.NewCommonHandler(h.container, h.countryUsecase)
}
