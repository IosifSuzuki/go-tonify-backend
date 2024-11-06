package controller

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/repository"
	"net/http"
)

type Account struct {
	container         container.Container
	accountRepository repository.AccountRepository
}

func NewAccount(container container.Container, accountRepository repository.AccountRepository) *Account {
	return &Account{
		container:         container,
		accountRepository: accountRepository,
	}
}

// GetMyAccount godoc
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
func (p *Account) GetMyAccount(ctx *gin.Context) {
	authToken, exist := ctx.Get(model.AuthorizationTokenKey)
	if !exist {
		err := model.MissedAuthorizationTokenError
		sendError(ctx, err, http.StatusUnauthorized)
		return
	}
	accessClaimsToken := authToken.(*model.AccessClaimsToken)
	domainAccount, err := p.accountRepository.FetchByID(ctx, accessClaimsToken.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	account := model.Account{
		ID:         domainAccount.ID,
		TelegramID: domainAccount.TelegramID,
		FirstName:  *domainAccount.FirstName,
		MiddleName: domainAccount.MiddleName,
		LastName:   *domainAccount.LastName,
		Nickname:   domainAccount.Nickname,
		AboutMe:    domainAccount.AboutMe,
		Gender:     model.NewGender(*domainAccount.Gender),
		Country:    domainAccount.Country,
		Location:   domainAccount.Location,
		CompanyID:  domainAccount.CompanyID,
	}
	sendResponse(ctx, account)
}

// EditMyAccount godoc
//
//	@Summary		edit my account
//	@Description	get actual account model
//	@Tags			profile
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.EditAccount	true	"updated account field"
//	@Success		200		{object}	model.Account
//	@Failure		500		"internal error"
//	@Failure		401		"invalid access token provided"
//	@Router			/account/edit [patch]
//	@Security		ApiKeyAuth
func (p *Account) EditMyAccount(ctx *gin.Context) {
	authToken, exist := ctx.Get(model.AuthorizationTokenKey)
	if !exist {
		err := model.MissedAuthorizationTokenError
		sendError(ctx, err, http.StatusUnauthorized)
		return
	}
	accessClaimsToken := authToken.(*model.AccessClaimsToken)
	var editAccount model.EditAccount
	if err := ctx.ShouldBindJSON(&editAccount); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	gender := editAccount.Gender.String()
	domainAccount := domain.Account{
		ID:         &accessClaimsToken.ID,
		FirstName:  &editAccount.FirstName,
		MiddleName: editAccount.MiddleName,
		LastName:   &editAccount.LastName,
		Nickname:   editAccount.Nickname,
		AboutMe:    editAccount.AboutMe,
		Gender:     &gender,
		Country:    &editAccount.Country,
		Location:   &editAccount.Location,
	}
	if err := p.accountRepository.UpdateAccount(ctx, &domainAccount); err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	updatedDomainAccount, err := p.accountRepository.FetchByID(ctx, accessClaimsToken.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	account := model.Account{
		ID:         updatedDomainAccount.ID,
		TelegramID: updatedDomainAccount.TelegramID,
		FirstName:  *updatedDomainAccount.FirstName,
		MiddleName: updatedDomainAccount.MiddleName,
		LastName:   *updatedDomainAccount.LastName,
		Nickname:   updatedDomainAccount.Nickname,
		AboutMe:    updatedDomainAccount.AboutMe,
		Gender:     model.NewGender(*updatedDomainAccount.Gender),
		Country:    updatedDomainAccount.Country,
		Location:   updatedDomainAccount.Location,
		CompanyID:  updatedDomainAccount.CompanyID,
	}
	sendResponse(ctx, account)
}
