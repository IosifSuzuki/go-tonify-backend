package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/account/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type AuthHandler struct {
	container      container.Container
	accountUsecase usecase.Account
}

func NewAuthHandler(
	container container.Container,
	accountUsecase usecase.Account,
) *AuthHandler {
	return &AuthHandler{
		container:      container,
		accountUsecase: accountUsecase,
	}
}

// SignUp godoc
//
//	@Summary		sign up / create account
//	@Description	record a account to db then return pairs of jwt tokens (access / refresh)
//	@Tags			auth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			telegram_raw_init_data	formData	string			true	"telegram_raw_init_data"
//	@Param			first_name				formData	string			true	"first_name"
//	@Param			middle_name				formData	string			false	"middle_name"
//	@Param			last_name				formData	string			true	"last_name"
//	@Param			role					formData	string			true	"role"
//	@Param			nickname				formData	string			true	"nickname"
//	@Param			about_me				formData	string			true	"about_me"
//	@Param			gender					formData	string			true	"gender"
//	@Param			country					formData	string			true	"country"
//	@Param			location				formData	string			true	"location"
//	@Param			company_name			formData	string			false	"company_name"
//	@Param			company_description		formData	string			false	"company_description"
//	@Param			avatar					formData	file			true	"avatar"
//	@Param			document				formData	file			true	"document"
//	@Success		201						{object}	dto.PairToken	"pair token"
//	@Failure		400						"bad parameters"
//	@Failure		500						"internal error"
//	@Router			/auth/sign-up [post]
func (a *AuthHandler) SignUp(ctx *gin.Context) {
	log := a.container.GetLogger()
	var createAccountRequest dto.CreateAccount
	if err := ctx.ShouldBind(&createAccountRequest); err != nil {
		log.Error("fail to parse/validate request model", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	avatarFileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		log.Error("fail to retrieve avatar file header", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	documentFileHeader, err := ctx.FormFile("document")
	if err != nil {
		log.Error("fail to retrieve document file header", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	createAccount := model.CreateAccount{
		TelegramInitData:   createAccountRequest.TelegramRawInitData,
		FirstName:          createAccountRequest.FirstName,
		MiddleName:         createAccountRequest.MiddleName,
		LastName:           createAccountRequest.LastName,
		Role:               string(createAccountRequest.Role),
		Nickname:           createAccountRequest.Nickname,
		AboutMe:            createAccountRequest.AboutMe,
		Gender:             string(createAccountRequest.Gender),
		Country:            createAccountRequest.Country,
		Location:           createAccountRequest.Location,
		CompanyName:        createAccountRequest.CompanyName,
		CompanyDescription: createAccountRequest.CompanyDescription,
		AvatarFileHeader:   avatarFileHeader,
		DocumentFileHeader: documentFileHeader,
	}
	accountID, err := a.accountUsecase.CreateAccount(ctx, createAccount)
	if err != nil {
		switch err {
		case model.NilError:
			failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		case model.DuplicateAccountWithTelegramIDError:
			failResponse(ctx, http.StatusConflict, dto.DuplicateAccountWithTelegramIDError)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError)
		}
		return
	}
	if accountID == nil {
		log.Error("account id has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	pairTokenModel, err := a.accountUsecase.GeneratePairToken(*accountID)
	if err != nil {
		log.Error("fail to generate pair token", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	if pairTokenModel == nil {
		log.Error("pair token has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	pairToken := converter.ConvertModel2PairTokenResponse(*pairTokenModel)
	successResponse(ctx, http.StatusCreated, pairToken)
}

// SignIn godoc
//
//	@Summary		sign in & authentication
//	@Description	process of authentication to system through provided credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.Credential	true	"credential"
//	@Success		200		{object}	dto.PairToken	"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		401		"incorrect or missing credentials"
//	@Failure		500		"internal error"
//	@Router			/auth/sign-in [post]
func (a *AuthHandler) SignIn(ctx *gin.Context) {
	log := a.container.GetLogger()
	var credentialRequest dto.Credential
	if err := ctx.ShouldBindJSON(&credentialRequest); err != nil {
		log.Error("fail to parse/validate request model", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	accountID, err := a.accountUsecase.AuthenticationTelegram(ctx, credentialRequest.TelegramRawInitData)
	if err != nil {
		log.Error("fail to authentication account by telegram init data", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	if accountID == nil {
		log.Error("account id has nil value", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	pairTokenModel, err := a.accountUsecase.GeneratePairToken(*accountID)
	if err != nil {
		log.Error("fail to generate pair token", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	if pairTokenModel == nil {
		log.Error("pair token has nil value")
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	pairToken := converter.ConvertModel2PairTokenResponse(*pairTokenModel)
	successResponse(ctx, http.StatusOK, pairToken)
}
