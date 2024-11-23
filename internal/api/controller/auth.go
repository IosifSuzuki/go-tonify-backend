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
	container            container.Container
	authService          service.AuthService
	attachmentService    service.AttachmentService
	attachmentRepository repository.AttachmentRepository
	companyRepository    repository.CompanyRepository
}

func NewAuthController(
	container container.Container,
	authService service.AuthService,
	attachmentService service.AttachmentService,
	attachmentRepository repository.AttachmentRepository,
	companyRepository repository.CompanyRepository,
) *AuthController {
	return &AuthController{
		container:            container,
		authService:          authService,
		attachmentService:    attachmentService,
		attachmentRepository: attachmentRepository,
		companyRepository:    companyRepository,
	}
}

// AccountSignUp godoc
//
//	@Summary		account sign up
//	@Description	record account to db and return pairs jwt tokens
//	@Tags			auth
//	@Accept			multipart/form-data
//	@Produce		json
//
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
//
//	@Param			avatar					formData	file			true	"Profile picture"
//	@Param			document				formData	file			true	"Profile Document"
//	@Success		201						{object}	model.PairToken	"pair token"
//	@Failure		400						"bad parameters"
//	@Failure		500						"internal error"
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
	documentFileHeader, err := ctx.FormFile("document")
	if err != nil {
		log.Error("fail to find document", logger.FError(err))
		sendError(ctx, err, http.StatusBadRequest)
		return
	}
	avatarFileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		log.Error("fail to find avatar", logger.FError(err))
		sendError(ctx, err, http.StatusBadRequest)
		return
	}
	createAccountRequest.Avatar = avatarFileHeader
	createAccountRequest.Document = documentFileHeader
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

	documentAttachmentDomain := prepareAttachment(createAccountRequest.Document.Filename)
	documentID, err := a.attachmentRepository.Create(ctx, documentAttachmentDomain)
	if err != nil {
		log.Error("fail to create document attachment in db ", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentAttachmentDomain.ID = *documentID

	avatarAttachmentDomain := prepareAttachment(createAccountRequest.Avatar.Filename)
	avatarID, err := a.attachmentRepository.Create(ctx, documentAttachmentDomain)
	if err != nil {
		log.Error("fail to create avatar attachment in db ", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	avatarAttachmentDomain.ID = *avatarID
	var createAccount = model.CreateAccount{
		TelegramRawInitData: createAccountRequest.TelegramRawInitData,
		FirstName:           createAccountRequest.FirstName,
		MiddleName:          createAccountRequest.MiddleName,
		LastName:            createAccountRequest.LastName,
		Nickname:            createAccountRequest.Nickname,
		AboutMe:             createAccountRequest.AboutMe,
		Role:                createAccountRequest.Role,
		Gender:              createAccountRequest.Gender,
		Country:             createAccountRequest.Country,
		Location:            createAccountRequest.Location,
		CompanyID:           companyID,
		AvatarID:            avatarID,
		DocumentID:          documentID,
	}
	accountID, err := a.authService.CreateAccount(context.Background(), &createAccount)
	if err != nil {
		log.Error("fail to create account", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}

	avatarUploadFile := model.UploadFile{
		Name: avatarAttachmentDomain.FileName,
		Body: avatarFile,
	}
	remoteAvatarURL, err := a.attachmentService.UploadFile(avatarUploadFile)
	if err != nil {
		log.Error("fail to upload avatar file", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	avatarAttachmentDomain.Path = remoteAvatarURL
	avatarAttachmentDomain.Status = string(model.InUseAttachmentStatus)
	if err := a.attachmentRepository.Update(ctx, avatarAttachmentDomain); err != nil {
		log.Error("fail to update avatar attachment", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentUploadFile := model.UploadFile{
		Name: documentAttachmentDomain.FileName,
		Body: documentFile,
	}
	documentAvatarURL, err := a.attachmentService.UploadFile(documentUploadFile)
	if err != nil {
		log.Error("fail to upload document file", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	documentAttachmentDomain.Path = documentAvatarURL
	documentAttachmentDomain.Status = string(model.InUseAttachmentStatus)
	if err := a.attachmentRepository.Update(ctx, documentAttachmentDomain); err != nil {
		log.Error("fail to update attachment", logger.FError(err))
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.authService.GenerateAccountJWT(context.Background(), *accountID)
	if err != nil {
		log.Error("fail to generate account jwt", logger.FError(err))
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
	pairToken, err := a.authService.GenerateAccountJWT(context.Background(), account.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusOK)
}
