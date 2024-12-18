package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/account/usecase"
	"go-tonify-backend/pkg/logger"
	"mime/multipart"
	"net/http"
)

type AuthHandler struct {
	container      container.Container
	validation     validator.HttpValidator
	accountUsecase usecase.Account
}

func NewAuthHandler(
	container container.Container,
	validation validator.HttpValidator,
	accountUsecase usecase.Account,
) *AuthHandler {
	return &AuthHandler{
		container:      container,
		validation:     validation,
		accountUsecase: accountUsecase,
	}
}

// SignUp godoc
//
//	@Summary		Sign up a new user
//	@Description	Create a new account. The account will be validated and recorded in the database.
//	@Description	Attachments will be saved in external storage. Server returns access / refresh tokens
//	@Tags			auth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			telegram_init_data	formData	string									true	"telegram init data"
//	@Param			first_name			formData	string									true	"first name"
//	@Param			middle_name			formData	string									false	"middle name"
//	@Param			last_name			formData	string									true	"last mame"
//	@Param			role				formData	string									true	"role"	Enums(client, freelancer)
//	@Param			nickname			formData	string									true	"nickname"
//	@Param			about_me			formData	string									false	"about me"
//	@Param			gender				formData	string									true	"gender"	Enums(male, female)
//	@Param			country				formData	string									true	"country"
//	@Param			location			formData	string									true	"location"
//	@Param			tags				formData	[]string								false	"tags"
//	@Param			company_name		formData	string									false	"company name"
//	@Param			company_description	formData	string									false	"company description"
//	@Param			avatar				formData	file									false	"avatar file"
//	@Param			document			formData	file									false	"document file"
//	@Success		201					{object}	dto.Response{response=dto.PairToken}	"access & refresh tokens"
//	@Failure		400					{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		409					{object}	dto.Response{response=dto.Empty}		"detailed error message, provided data already exist"
//	@Failure		500					{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/auth/sign-up [post]
func (a *AuthHandler) SignUp(ctx *gin.Context) {
	log := a.container.GetLogger()
	var (
		createAccountRequest dto.CreateAccount
		avatarFileHeader     *multipart.FileHeader
		documentFileHeader   *multipart.FileHeader
	)
	err := ctx.ShouldBind(&createAccountRequest)
	if err != nil {
		log.Error("fail to parse/validate request model", logger.FError(err))
		badRequestResponse(ctx, a.validation, dto.BadRequestError, err)
		return
	}
	if err := ctx.Request.ParseMultipartForm(50 << 20); err != nil {
		log.Error("fail to Parse multipart form", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError, err)
	}
	if _, ok := ctx.Request.MultipartForm.File["avatar"]; ok {
		log.Debug("has avatar")
		avatarFileHeader, err = ctx.FormFile("avatar")
		if err != nil {
			log.Error("fail to retrieve avatar file header", logger.FError(err))
			failResponse(ctx, http.StatusBadRequest, dto.BadRequestError, err)
			return
		}
	}
	if _, ok := ctx.Request.MultipartForm.File["document"]; ok {
		documentFileHeader, err = ctx.FormFile("document")
		if err != nil {
			log.Error("fail to retrieve document file header", logger.FError(err))
			failResponse(ctx, http.StatusBadRequest, dto.BadRequestError, err)
			return
		}
	}
	createAccount := model.CreateAccount{
		TelegramInitData:   createAccountRequest.TelegramInitData,
		FirstName:          createAccountRequest.FirstName,
		MiddleName:         createAccountRequest.MiddleName,
		LastName:           createAccountRequest.LastName,
		Role:               string(createAccountRequest.Role),
		Nickname:           &createAccountRequest.Nickname,
		AboutMe:            createAccountRequest.AboutMe,
		Gender:             string(createAccountRequest.Gender),
		Country:            createAccountRequest.Country,
		Location:           createAccountRequest.Location,
		Tags:               createAccountRequest.Tags,
		CompanyName:        createAccountRequest.CompanyName,
		CompanyDescription: createAccountRequest.CompanyDescription,
		AvatarFileHeader:   avatarFileHeader,
		DocumentFileHeader: documentFileHeader,
	}
	accountID, err := a.accountUsecase.CreateAccount(ctx, createAccount)
	if err != nil {
		switch err {
		case model.NilError:
			failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		case model.DuplicateAccountWithTelegramIDError:
			failResponse(ctx, http.StatusConflict, dto.DuplicateAccountWithTelegramIDError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	if accountID == nil {
		log.Error("account id has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	pairTokenModel, err := a.accountUsecase.GeneratePairToken(*accountID)
	if err != nil {
		log.Error("fail to generate pair token", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	if pairTokenModel == nil {
		log.Error("pair token has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	pairToken := converter.ConvertModel2PairTokenResponse(*pairTokenModel)
	successResponse(ctx, http.StatusCreated, pairToken)
}

// SignIn godoc
//
//	@Summary		Sign in an existing user
//	@Description	The user provides their Telegram initialization data. The user performs authentication, and if successful, access and refresh tokens are returned
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.Credential							true	"credential"
//	@Success		200		{object}	dto.Response{response=dto.PairToken}	"pair token"
//	@Failure		400		{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		410		{object}	dto.Response{response=dto.Empty}	"account does not exist or has been deleted"
//	@Failure		500		{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/auth/sign-in [post]
func (a *AuthHandler) SignIn(ctx *gin.Context) {
	log := a.container.GetLogger()
	var credentialRequest dto.Credential
	if err := ctx.ShouldBindJSON(&credentialRequest); err != nil {
		log.Error("fail to parse/validate request model", logger.FError(err))
		badRequestResponse(ctx, a.validation, dto.BadRequestError, err)
		return
	}
	accountID, err := a.accountUsecase.AuthenticationTelegram(ctx, credentialRequest.TelegramInitData)
	if err != nil {
		log.Error("fail to authentication account by telegram init data", logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	if accountID == nil {
		log.Error("account id has nil value", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	pairTokenModel, err := a.accountUsecase.GeneratePairToken(*accountID)
	if err != nil {
		log.Error("fail to generate pair token", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	if pairTokenModel == nil {
		log.Error("pair token has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	pairToken := converter.ConvertModel2PairTokenResponse(*pairTokenModel)
	successResponse(ctx, http.StatusOK, pairToken)
}
