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

type AccountHandler struct {
	container      container.Container
	accountUsecase usecase.Account
}

func NewAccountHandler(
	container container.Container,
	accountUsecase usecase.Account,
) *AccountHandler {
	return &AccountHandler{
		container:      container,
		accountUsecase: accountUsecase,
	}
}

// GetMy godoc
//
//	@Summary		get my account
//	@Description	get actual account model
//	@Tags			profile
//	@Produce		json
//	@Success		200	{object}	model.Account
//	@Failure		500	"internal error"
//	@Failure		401	"invalid access token provided"
//	@Router			/account/my [get]
//	@Security		ApiKeyAuth
func (a *AccountHandler) GetMy(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err)
		return
	}
	accountModel, err := a.accountUsecase.GetDetailsAccount(ctx, *accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.F("account_id", accountID), logger.FError(err))
		return
	}
	account := converter.ConvertModel2AccountResponse(accountModel)
	successResponse(ctx, http.StatusOK, account)
}

// EditMy godoc
//
//	@Summary		edit my account
//	@Description	get actual account model
//	@Tags			profile
//	@Accept			multipart/form-data
//	@Produce		json
//
//	@Param			first_name			formData	string	true	"first_name"
//	@Param			middle_name			formData	string	false	"middle_name"
//	@Param			last_name			formData	string	true	"last_name"
//	@Param			role				formData	string	true	"role"
//	@Param			nickname			formData	string	true	"nickname"
//	@Param			about_me			formData	string	true	"about_me"
//	@Param			gender				formData	string	true	"gender"
//	@Param			country				formData	string	true	"country"
//	@Param			location			formData	string	true	"location"
//	@Param			company_name		formData	string	true	"company_name"
//	@Param			company_description	formData	string	true	"company_description"
//
//	@Param			avatar				formData	file	true	"Profile picture"
//	@Param			document			formData	file	true	"Profile Document"
//	@Success		200					{object}	model.Account
//	@Failure		500					"internal error"
//	@Failure		401					"invalid access token provided"
//	@Router			/account/edit [patch]
//	@Security		ApiKeyAuth
func (a *AccountHandler) EditMy(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err)
		return
	}
	avatarFileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		log.Error("can't retrieve a avatar file from form", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	documentFileHeader, err := ctx.FormFile("document")
	if err != nil {
		log.Error("can't retrieve a document file from form", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	var editAccountRequest dto.EditAccount
	if err := ctx.ShouldBind(&editAccountRequest); err != nil {
		log.Error("fail to bind edit account", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	var editAccount = model.EditAccount{
		ID:                 *accountID,
		FirstName:          editAccountRequest.FirstName,
		MiddleName:         editAccountRequest.MiddleName,
		LastName:           editAccountRequest.LastName,
		Role:               string(editAccountRequest.Role),
		Nickname:           editAccountRequest.Nickname,
		AboutMe:            editAccountRequest.AboutMe,
		Gender:             string(editAccountRequest.Gender),
		Location:           editAccountRequest.Location,
		Country:            editAccountRequest.Country,
		CompanyName:        editAccountRequest.CompanyName,
		CompanyDescription: editAccountRequest.CompanyDescription,
		AvatarFileHeader:   avatarFileHeader,
		DocumentFileHeader: documentFileHeader,
	}
	err = a.accountUsecase.EditAccount(ctx, editAccount)
	if err != nil {
		log.Error("fail process edit account", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError)
		return
	}
	accountModel, err := a.accountUsecase.GetDetailsAccount(ctx, *accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.F("account_id", *accountID), logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.CastTypeError)
		return
	}
	account := converter.ConvertModel2AccountResponse(accountModel)
	successResponse(ctx, http.StatusOK, account)
}

// MatchAccounts godoc
//
//	@Summary		match accounts
//	@Description	get accounts by matching
//	@Tags			profile
//	@Param			request	body	model.MatchAccountRequest	true	"match account"
//	@Produce		json
//	@Success		200	{array}	model.Account
//	@Failure		500	"internal error"
//	@Failure		401	"invalid access token provided"
//	@Router			/account/matching [get]
//	@Security		ApiKeyAuth
func (a *AccountHandler) MatchAccounts(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err)
		return
	}
	var getMatchAccounts dto.GetMatchAccounts
	if err := ctx.ShouldBindJSON(&getMatchAccounts); err != nil {
		log.Error("fail to bind get match accounts", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, dto.BadRequestError)
		return
	}
	accountModels, err := a.accountUsecase.GetMatchAccounts(ctx, *accountID, getMatchAccounts.Limit)
	if err != nil {
		log.Error("fail to get match accounts", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError)
		return
	}
	accounts := make([]dto.Account, 0, len(accountModels))
	for _, accountModel := range accountModels {
		account := converter.ConvertModel2AccountResponse(&accountModel)
		accounts = append(accounts, *account)
	}
	successResponse(ctx, http.StatusOK, accounts)
}
