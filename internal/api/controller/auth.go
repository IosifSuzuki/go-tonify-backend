package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type AuthController struct {
	container         container.Container
	authService       service.AuthService
	attachmentService service.AttachmentService
	companyRepository repository.CompanyRepository
}

func NewAuthController(
	container container.Container,
	authService service.AuthService,
	attachmentService service.AttachmentService,
	companyRepository repository.CompanyRepository,
) *AuthController {
	return &AuthController{
		container:         container,
		authService:       authService,
		attachmentService: attachmentService,
		companyRepository: companyRepository,
	}
}

// AccountSignUp godoc
//
//	@Summary		account sign up
//	@Description	record account to db and return pairs jwt tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateAccountRequest	true	"account payload"
//	@Success		201		{object}	model.PairToken				"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		500		"internal error"
//	@Router			/auth/sign-up [post]
func (a *AuthController) AccountSignUp(ctx *gin.Context) {
	log := a.container.GetLogger()
	var createAccountRequest model.CreateAccountRequest
	if err := ctx.ShouldBind(&createAccountRequest); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}

	var (
		companyID *int64
		err       error
	)
	if createAccountRequest.CompanyName != nil && createAccountRequest.CompanyDescription != nil {
		companyEntity := domain.Company{
			Name:        createAccountRequest.CompanyName,
			Description: createAccountRequest.CompanyDescription,
		}
		companyID, err = a.companyRepository.Create(ctx, &companyEntity)
		if err != nil {
			log.Error("fail to record company to db", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
	}

	avatarFile, err := createAccountRequest.Avatar.Open()
	if err != nil {
		log.Error("fail to open avatar file from request", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentFile, err := createAccountRequest.Avatar.Open()
	if err != nil {
		log.Error("fail to open document file from request", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	avatarUploadFile := model.UploadFile{
		Name: createAccountRequest.Avatar.Filename,
		Body: avatarFile,
	}
	documentUploadFile := model.UploadFile{
		Name: createAccountRequest.Document.Filename,
		Body: documentFile,
	}
	avatarURL, err := a.attachmentService.UploadFile(avatarUploadFile)
	if err != nil {
		log.Error("fail to upload avatar file", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentURL, err := a.attachmentService.UploadFile(documentUploadFile)
	if err != nil {
		log.Error("fail to upload document file", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}

	var createAccount = model.CreateAccount{
		TelegramRawInitData: createAccountRequest.TelegramRawInitData,
		FirstName:           createAccountRequest.FirstName,
		MiddleName:          createAccountRequest.MiddleName,
		LastName:            createAccountRequest.LastName,
		Nickname:            createAccountRequest.Nickname,
		AboutMe:             createAccountRequest.AboutMe,
		Gender:              createAccountRequest.Gender,
		Country:             createAccountRequest.Country,
		Location:            createAccountRequest.Location,
		CompanyID:           companyID,
		AvatarURL:           avatarURL,
		DocumentURL:         documentURL,
	}
	accountID, err := a.authService.CreateAccount(context.Background(), &createAccount)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.authService.GenerateAccountJWT(context.Background(), *accountID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusCreated)
}

// AccountSignIn godoc
//
//	@Summary		account sign in
//	@Description	process of authorization to system through provided credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.Credential	true	"credential"
//	@Success		200		{object}	model.PairToken		"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		401		"incorrect or missing credentials"
//	@Failure		500		"internal error"
//	@Router			/auth/sign-in [post]
func (a *AuthController) AccountSignIn(ctx *gin.Context) {
	var credential model.Credential
	if err := ctx.ShouldBindJSON(&credential); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	account, err := a.authService.AuthorizationAccount(ctx, &credential)
	if err == model.AccountNotExistsError {
		sendError(ctx, err, http.StatusUnauthorized)
		return
	} else if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.authService.GenerateAccountJWT(context.Background(), *account.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusOK)
}
