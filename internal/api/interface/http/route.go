package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-tonify-backend/docs"
	"go-tonify-backend/internal/api/interface/http/middleware"
	"go-tonify-backend/internal/api/interface/http/v1"
	"go-tonify-backend/internal/container"
	accountUsecase "go-tonify-backend/internal/domain/account/usecase"
	countryUsecase "go-tonify-backend/internal/domain/country/usecase"
	"go-tonify-backend/pkg/datetime"
	"time"
)

//	@title			Swagger Tonify API
//	@version		1.0
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	tonifyapp@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

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

	h.ConfigureSwagDocs()

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

func (h *Handler) ConfigureSwagDocs() {
	serverConf := h.container.GetServerConfig()
	timeStarted := datetime.GetTimeString(time.Now(), datetime.ReadableFormatLayoutTime)
	docs.SwaggerInfo.Host = serverConf.Addr
	docs.SwaggerInfo.Description = fmt.Sprintf("The server helps interact with the telegram mini-application. API supports versions: /v1/, /v2/\nLast updated: %s", timeStarted)
	docs.SwaggerInfo.BasePath = "/api/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
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
