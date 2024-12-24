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

type AccountHandler struct {
	container      container.Container
	validation     validator.HttpValidator
	accountUsecase usecase.Account
}

func NewAccountHandler(
	container container.Container,
	validation validator.HttpValidator,
	accountUsecase usecase.Account,
) *AccountHandler {
	return &AccountHandler{
		container:      container,
		validation:     validation,
		accountUsecase: accountUsecase,
	}
}

// GetMy godoc
//
//	@Summary		Get My Account
//	@Description	Get the details of the authenticated user's account
//	@Tags			account
//	@Produce		json
//	@Param			Authorization	header		string					true	"account's access token"
//	@Success		200	{object}	dto.Response{response=dto.Account}	"account details"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}	"the authorization token is invalid/expired/missing"
//	@Failure		410	{object}	dto.Response{response=dto.Empty}	"account does not exist or has been deleted"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Router			/v1/account/my [get]
//	@Security		ApiKeyAuth
func (a *AccountHandler) GetMy(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err, nil)
		return
	}
	accountModel, err := a.accountUsecase.GetDetailsAccount(ctx, *accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.F("account_id", accountID), logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	account := converter.ConvertModel2AccountResponse(accountModel)
	successResponse(ctx, http.StatusOK, account)
}

// EditMy godoc
//
//	@Summary		Edit my account
//	@Description	Edit the details of the authenticated user's account
//	@Tags			account
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization		header		string								true	"account's access token"
//	@Param			first_name			formData	string								true	"first name"
//	@Param			middle_name			formData	string								false	"middle name"
//	@Param			last_name			formData	string								true	"last mame"
//	@Param			role				formData	string								true	"role"	Enums(client, freelancer)
//	@Param			nickname			formData	string								true	"nickname"
//	@Param			about_me			formData	string								false	"about me"
//	@Param			gender				formData	string								true	"gender"	Enums(male, female)
//	@Param			country				formData	string								true	"country"
//	@Param			location			formData	string								true	"location"
//	@Param			company_name		formData	string								false	"company name"
//	@Param			company_description	formData	string								false	"company description"
//	@Param			avatar				formData	file								true	"avatar file"
//	@Param			document			formData	file								true	"document file"
//	@Success		200					{object}	dto.Response{response=dto.Account}	"account details"
//	@Failure		400					{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Failure		401					{object}	dto.Response{response=dto.Empty}	"the authorization token is invalid/expired/missing"
//	@Failure		410					{object}	dto.Response{response=dto.Empty}	"account does not exist or has been deleted"
//	@Failure		500					{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Router			/v1/account/edit [patch]
//	@Security		ApiKeyAuth
func (a *AccountHandler) EditMy(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err, nil)
		return
	}
	var (
		avatarFileHeader   *multipart.FileHeader
		documentFileHeader *multipart.FileHeader
	)
	if _, ok := ctx.Request.MultipartForm.File["avatar"]; ok {
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
	var editAccountRequest dto.EditAccount
	if err := ctx.ShouldBind(&editAccountRequest); err != nil {
		log.Error("fail to bind edit account", logger.FError(err))
		badRequestResponse(ctx, a.validation, dto.BadRequestError, err)
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
		Tags:               editAccountRequest.Tags,
		CategoryIDs:        editAccountRequest.CategoryIDs,
		AvatarFileHeader:   avatarFileHeader,
		DocumentFileHeader: documentFileHeader,
	}
	err = a.accountUsecase.EditAccount(ctx, editAccount)
	if err != nil {
		log.Error("fail process edit account", logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	accountModel, err := a.accountUsecase.GetDetailsAccount(ctx, *accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.F("account_id", *accountID), logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
	}
	account := converter.ConvertModel2AccountResponse(accountModel)
	successResponse(ctx, http.StatusOK, account)
}

// ChangeRole godoc
//
//	@Summary		Change role for my account
//	@Description	A user will change their own role
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"account's access token"
//	@Param			request			body		dto.ChangeRole			true	"parameter with new role"
//	@Success		200	{object}	dto.Response{response=dto.Account}	"updated account details"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}	"the authorization token is invalid/expired/missing"
//	@Failure		410	{object}	dto.Response{response=dto.Empty}	"account does not exist or has been deleted"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Router			/v1/account/change/role [patch]
//	@Security		ApiKeyAuth
func (a *AccountHandler) ChangeRole(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err, nil)
		return
	}
	var changeRole dto.ChangeRole
	if err := ctx.ShouldBind(&changeRole); err != nil {
		log.Error("fail to bind request model", logger.FError(err))
		badRequestResponse(ctx, a.validation, dto.BadRequestError, err)
		return
	}
	role := converter.ConvertDto2RoleModel(changeRole.NewRole)
	if err := a.accountUsecase.ChangeRole(ctx, *accountID, role); err != nil {
		log.Error(
			"fail to change role for account",
			logger.F("account_id", *accountID),
			logger.FError(err),
		)
		failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		return
	}
	accountModel, err := a.accountUsecase.GetDetailsAccount(ctx, *accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.F("account_id", accountID), logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	account := converter.ConvertModel2AccountResponse(accountModel)
	successResponse(ctx, http.StatusOK, account)
}

// DeleteMy godoc
//
//	@Summary		Delete my account
//	@Description	Delete the authenticated user's account
//	@Param			Authorization	header	string	true	"account's access token"
//	@Tags			account
//	@Produce		json
//	@Success		200					{object}	dto.Response{response=string}		"returns ok string"
//	@Failure		400					{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Failure		401					{object}	dto.Response{response=dto.Empty}	"the authorization token is invalid/expired/missing"
//	@Failure		410					{object}	dto.Response{response=dto.Empty}	"account does not exist or has been deleted"
//	@Failure		500					{object}	dto.Response{response=dto.Empty}	"detailed error message"
//	@Router			/v1/account/delete [delete]
//	@Security		ApiKeyAuth
func (a *AccountHandler) DeleteMy(ctx *gin.Context) {
	log := a.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusBadRequest, err, nil)
		return
	}
	if err := a.accountUsecase.DeleteAccount(ctx, *accountID); err != nil {
		log.Error("fail to process delete account", logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	successResponse(ctx, http.StatusOK, "ok")
}
