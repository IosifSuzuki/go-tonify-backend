package controller

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/service"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type Account struct {
	container            container.Container
	accountRepository    repository.AccountRepository
	attachmentRepository repository.AttachmentRepository
	companyRepository    repository.CompanyRepository
	attachmentService    service.AttachmentService
}

func NewAccount(container container.Container, accountRepository repository.AccountRepository, attachmentRepository repository.AttachmentRepository, companyRepository repository.CompanyRepository, attachmentService service.AttachmentService) *Account {
	return &Account{
		container:            container,
		accountRepository:    accountRepository,
		attachmentRepository: attachmentRepository,
		companyRepository:    companyRepository,
		attachmentService:    attachmentService,
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
	log := p.container.GetLogger()
	authToken, exist := ctx.Get(model.AuthorizationTokenKey)
	if !exist {
		err := model.MissedAuthorizationTokenError
		sendError(ctx, err, http.StatusUnauthorized)
		return
	}
	accessClaimsToken := authToken.(*model.AccessClaimsToken)
	accountDomain, err := p.accountRepository.FetchByID(ctx, accessClaimsToken.ID)
	if err != nil {
		log.Error("fail to fetch account by id", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	if accountDomain.AvatarAttachmentID == nil || accountDomain.DocumentAttachmentID == nil {
		err := model.NilError
		log.Error("bad state for occurred")
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	companyDomain, err := p.companyRepository.FetchByID(ctx, *accountDomain.CompanyID)
	if err != nil {
		log.Error("fail to fetch company by id", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	avatarAttachment, err := p.attachmentRepository.FetchByID(ctx, *accountDomain.AvatarAttachmentID)
	if err != nil {
		log.Error("fail to fetch attachment", logger.FError(err), logger.F("id", *accountDomain.AvatarAttachmentID))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentAttachment, err := p.attachmentRepository.FetchByID(ctx, *accountDomain.DocumentAttachmentID)
	if err != nil {
		log.Error("fail to fetch attachment", logger.FError(err), logger.F("id", *accountDomain.AvatarAttachmentID))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	account := model.Account{
		ID:                 accountDomain.ID,
		TelegramID:         accountDomain.TelegramID,
		FirstName:          accountDomain.FirstName,
		MiddleName:         accountDomain.MiddleName,
		LastName:           accountDomain.LastName,
		Nickname:           accountDomain.Nickname,
		Role:               accountDomain.Role,
		AboutMe:            accountDomain.AboutMe,
		Gender:             model.NewGender(accountDomain.Gender),
		Country:            accountDomain.Country,
		Location:           accountDomain.Location,
		CompanyID:          accountDomain.CompanyID,
		CompanyName:        companyDomain.Name,
		CompanyDescription: companyDomain.Description,
		AvatarURL:          avatarAttachment.Path,
		DocumentURL:        documentAttachment.Path,
	}
	sendResponse(ctx, account)
}

// EditMyAccount godoc
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
func (p *Account) EditMyAccount(ctx *gin.Context) {
	log := p.container.GetLogger()
	authToken, exist := ctx.Get(model.AuthorizationTokenKey)
	if !exist {
		err := model.MissedAuthorizationTokenError
		sendError(ctx, err, http.StatusUnauthorized)
		return
	}
	accessClaimsToken := authToken.(*model.AccessClaimsToken)
	accountDomain, err := p.accountRepository.FetchByID(ctx, accessClaimsToken.ID)
	if err != nil {
		log.Error("fetch account by id", logger.FError(err))
		sendError(ctx, model.DataBaseOperationError, http.StatusInternalServerError)
		return
	}
	var editAccountRequest model.EditAccountRequest
	if err := ctx.ShouldBind(&editAccountRequest); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	var (
		avatarRemoteURL   *string
		documentRemoteURL *string
	)
	log.Debug("edit profile receives", logger.F("edit_profile", editAccountRequest))
	if accountDomain.AvatarAttachmentID != nil {
		if err = p.attachmentRepository.Delete(ctx, *accountDomain.AvatarAttachmentID); err != nil {
			log.Error("delete attachment by id", logger.FError(err))
			sendError(ctx, model.DataBaseOperationError, http.StatusInternalServerError)
			return
		}
		avatarFile, err := editAccountRequest.Avatar.Open()
		if err != nil {
			log.Error("fail to open avatar file from request", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		avatarAttachmentDomain := prepareAttachment(editAccountRequest.Document.Filename)
		avatarAttachmentID, err := p.attachmentRepository.Create(ctx, avatarAttachmentDomain)
		if err != nil {
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		avatarAttachmentDomain.ID = *avatarAttachmentID
		avatarUploadFile := model.UploadFile{
			Name: avatarAttachmentDomain.FileName,
			Body: avatarFile,
		}
		if err != nil {
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		avatarRemoteURL, err = p.attachmentService.UploadFile(avatarUploadFile)
		if err != nil {
			log.Error("fail to upload avatar file", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		avatarAttachmentDomain.Path = avatarRemoteURL
		avatarAttachmentDomain.Status = string(model.InUseAttachmentStatus)
		if err := p.attachmentRepository.Update(ctx, avatarAttachmentDomain); err != nil {
			log.Error("fail to update avatar attachment", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		accountDomain.AvatarAttachmentID = avatarAttachmentID
	}
	if accountDomain.DocumentAttachmentID != nil {
		if err = p.attachmentRepository.Delete(ctx, *accountDomain.DocumentAttachmentID); err != nil {
			log.Error("delete attachment by id", logger.FError(err))
			sendError(ctx, model.DataBaseOperationError, http.StatusInternalServerError)
			return
		}
		documentFile, err := editAccountRequest.Avatar.Open()
		if err != nil {
			log.Error("fail to open document file from request", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		documentAttachmentDomain := prepareAttachment(editAccountRequest.Avatar.Filename)
		documentAttachmentID, err := p.attachmentRepository.Create(ctx, documentAttachmentDomain)
		if err != nil {
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		documentAttachmentDomain.ID = *documentAttachmentID
		documentUploadFile := model.UploadFile{
			Name: documentAttachmentDomain.FileName,
			Body: documentFile,
		}
		documentRemoteURL, err = p.attachmentService.UploadFile(documentUploadFile)
		if err != nil {
			log.Error("fail to upload document file", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		documentAttachmentDomain.Path = documentRemoteURL
		documentAttachmentDomain.Status = string(model.InUseAttachmentStatus)
		if err := p.attachmentRepository.Update(ctx, documentAttachmentDomain); err != nil {
			log.Error("fail to update document attachment", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		accountDomain.DocumentAttachmentID = documentAttachmentID
	}
	accountDomain.FirstName = editAccountRequest.FirstName
	accountDomain.MiddleName = editAccountRequest.MiddleName
	accountDomain.LastName = editAccountRequest.LastName
	accountDomain.Nickname = editAccountRequest.Nickname
	accountDomain.Role = string(editAccountRequest.Role)
	accountDomain.AboutMe = editAccountRequest.AboutMe
	accountDomain.Gender = editAccountRequest.Gender.String()
	accountDomain.Country = &editAccountRequest.Country
	accountDomain.Location = &editAccountRequest.Location
	var companyDomain *domain.Company
	if accountDomain.CompanyID != nil {
		companyDomain = &domain.Company{
			ID:          *accountDomain.CompanyID,
			Name:        editAccountRequest.CompanyName,
			Description: editAccountRequest.CompanyDescription,
		}
		log.Debug("will try to update company", logger.F("company", companyDomain))
		if err := p.companyRepository.Update(ctx, companyDomain); err != nil {
			log.Error("fail to update company", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
	} else if editAccountRequest.CompanyName != nil && editAccountRequest.CompanyDescription != nil {
		companyDomain = &domain.Company{
			Name:        editAccountRequest.CompanyName,
			Description: editAccountRequest.CompanyDescription,
		}
		log.Debug("will try to update company", logger.F("company", companyDomain))
		companyID, err := p.companyRepository.Create(ctx, companyDomain)
		if err != nil {
			log.Error("fail to create company", logger.FError(err))
			sendError(ctx, err, http.StatusInternalServerError)
			return
		}
		companyDomain.ID = *companyID
	} else {
		companyDomain = domain.NewCompany()
	}
	log.Debug("will try to update account", logger.F("account", accountDomain))
	if err := p.accountRepository.UpdateAccount(ctx, accountDomain); err != nil {
		log.Error("fail to update account", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	updatedDomainAccount, err := p.accountRepository.FetchByID(ctx, accessClaimsToken.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	account := model.Account{
		ID:                 updatedDomainAccount.ID,
		TelegramID:         updatedDomainAccount.TelegramID,
		FirstName:          updatedDomainAccount.FirstName,
		MiddleName:         updatedDomainAccount.MiddleName,
		LastName:           updatedDomainAccount.LastName,
		Nickname:           updatedDomainAccount.Nickname,
		Role:               updatedDomainAccount.Role,
		AboutMe:            updatedDomainAccount.AboutMe,
		Gender:             model.NewGender(updatedDomainAccount.Gender),
		Country:            updatedDomainAccount.Country,
		Location:           updatedDomainAccount.Location,
		CompanyID:          updatedDomainAccount.CompanyID,
		CompanyName:        companyDomain.Name,
		CompanyDescription: companyDomain.Description,
		AvatarURL:          avatarRemoteURL,
		DocumentURL:        documentRemoteURL,
	}
	sendResponse(ctx, account)
}
