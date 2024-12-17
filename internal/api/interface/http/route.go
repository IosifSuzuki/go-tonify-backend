package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	v "github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-tonify-backend/docs"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/middleware"
	"go-tonify-backend/internal/api/interface/http/v1"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/container"
	accountUsecase "go-tonify-backend/internal/domain/account/usecase"
	countryUsecase "go-tonify-backend/internal/domain/country/usecase"
	taskUsecase "go-tonify-backend/internal/domain/task/usecase"
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
	matchUsecase   accountUsecase.Match
	countryUsecase countryUsecase.Country
	taskUsecase    taskUsecase.Task
}

func NewHandler(
	container container.Container,
	accountUsecase accountUsecase.Account,
	matchUsecase accountUsecase.Match,
	countryUsecase countryUsecase.Country,
	taskUsecase taskUsecase.Task,
) *Handler {
	return &Handler{
		container:      container,
		accountUsecase: accountUsecase,
		matchUsecase:   matchUsecase,
		countryUsecase: countryUsecase,
		taskUsecase:    taskUsecase,
	}
}

func (h *Handler) Run() error {
	r := gin.Default()

	validation, err := h.configureAndInitValidation()
	if err != nil {
		return err
	}
	loggerMiddleware := middleware.NewLogger(h.container)
	corsMiddleware := middleware.NewCORS(h.container)
	authMiddleware := middleware.NewAuth(h.container, h.accountUsecase)
	roleMiddleware := middleware.NewRole(h.container, h.accountUsecase)

	r.Use(corsMiddleware.CORS())
	r.Use(loggerMiddleware.Logging())

	h.ConfigureSwagDocs()

	v1 := r.Group("api/v1")

	authHandler := h.composeAuthHandler(validation)

	authGroup := v1.Group("auth")
	{
		authGroup.POST("/sign-up", authHandler.SignUp)
		authGroup.POST("/sign-in", authHandler.SignIn)
	}
	accountHandler := h.composeAccount(validation)
	accountGroup := v1.Group("account")
	accountGroup.Use(authMiddleware.Authorization())
	{
		accountGroup.GET("/my", accountHandler.GetMy)
		accountGroup.PATCH("/edit", accountHandler.EditMy)
		accountGroup.PATCH("/change/role", accountHandler.ChangeRole)
		accountGroup.DELETE("/delete", accountHandler.DeleteMy)
	}
	matchHandler := h.composeMatch(validation)
	matchGroup := v1.Group("match")
	matchGroup.Use(authMiddleware.Authorization())
	{
		matchGroup.GET("/matchable/accounts", matchHandler.MatchableAccounts)
		matchGroup.POST("/action/:action", matchHandler.MatchAction)
	}
	taskHandler := h.composeTask(validation)
	taskGroup := v1.Group("task")
	taskGroup.Use(authMiddleware.Authorization())
	{
		taskGroup.POST("/create", roleMiddleware.Authorization(dto.ClientRole), taskHandler.CreateTask)
		taskGroup.GET("/list", taskHandler.GetListTask)
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

func (h *Handler) composeAuthHandler(validation validator.HttpValidator) *v1.AuthHandler {
	return v1.NewAuthHandler(h.container, validation, h.accountUsecase)
}

func (h *Handler) composeAccount(validation validator.HttpValidator) *v1.AccountHandler {
	return v1.NewAccountHandler(h.container, validation, h.accountUsecase)
}

func (h *Handler) composeCommon() *v1.CommonHandler {
	return v1.NewCommonHandler(h.container, h.countryUsecase)
}

func (h *Handler) composeMatch(validation validator.HttpValidator) *v1.MatchHandler {
	return v1.NewMatchHandler(h.container, validation, h.matchUsecase)
}

func (h *Handler) composeTask(validator validator.HttpValidator) *v1.TaskHandler {
	return v1.NewTaskHandler(h.container, validator, h.taskUsecase)
}

func (h *Handler) configureAndInitValidation() (validator.HttpValidator, error) {
	validationEngine, ok := binding.Validator.Engine().(*v.Validate)
	if !ok {
		return nil, dto.CastTypeError
	}
	validator := validator.NewValidator(h.container)
	if err := validator.Register(validationEngine); err != nil {
		return nil, err
	}
	return validator, nil
}
