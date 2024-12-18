package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/converter"
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/account/repository"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/internal/domain/filestorage"
	"go-tonify-backend/internal/domain/provider/transaction"
	"go-tonify-backend/internal/utils"
	"go-tonify-backend/pkg/jwt"
	"go-tonify-backend/pkg/logger"
	"go-tonify-backend/pkg/telegram"
	"mime/multipart"
	"strings"
)

type Account interface {
	CreateAccount(ctx context.Context, createAccount model.CreateAccount) (*int64, error)
	GeneratePairToken(accountID int64) (*model.PairToken, error)
	AuthenticationTelegram(ctx context.Context, telegramInitData string) (*int64, error)
	ParseAccessToken(accessToken string) (*int64, error)
	GetDetailsAccount(ctx context.Context, id int64) (*model.Account, error)
	EditAccount(ctx context.Context, editAccount model.EditAccount) error
	DeleteAccount(ctx context.Context, accountID int64) error
	AccountHasRole(ctx context.Context, accountID int64, role model.Role) (bool, error)
	ChangeRole(ctx context.Context, accountID int64, role model.Role) error
}

type account struct {
	container            container.Container
	fileStorage          filestorage.FileStorage
	accountRepository    repository.Account
	attachmentRepository repository.Attachment
	tagRepository        repository.Tag
	transactionProvider  *transaction.Provider
}

type uploadFile struct {
	Name string
	File multipart.File
}

func NewAccount(
	container container.Container,
	fileStorage filestorage.FileStorage,
	accountRepository repository.Account,
	attachmentRepository repository.Attachment,
	tagRepository repository.Tag,
	transactionProvider *transaction.Provider,
) Account {
	return &account{
		container:            container,
		fileStorage:          fileStorage,
		accountRepository:    accountRepository,
		attachmentRepository: attachmentRepository,
		tagRepository:        tagRepository,
		transactionProvider:  transactionProvider,
	}
}

func (a *account) CreateAccount(ctx context.Context, createAccount model.CreateAccount) (*int64, error) {
	log := a.container.GetLogger()
	var (
		accountID                *int64
		documentAttachmentEntity *entity.Attachment
		avatarAttachmentEntity   *entity.Attachment
	)
	var initData = telegram.InitData{
		Token: a.container.GetTelegramBotToken(),
	}
	telegramInitModel, err := initData.Decode(createAccount.TelegramInitData)
	if err != nil {
		log.Error("fail to decode telegram init data", logger.FError(err))
		return nil, model.DecodeTelegramInitDataError
	}
	validTelegramInitData, err := initData.Validate(telegramInitModel)
	if err != nil {
		log.Error("fail to validate telegram init data", logger.FError(err))
		return nil, model.InvalidTelegramInitDataError
	}
	if !validTelegramInitData {
		log.Error("not valid telegram init data", logger.FError(err))
		return nil, model.InvalidTelegramInitDataError
	}
	isDeletedAccountWithTelegramID, err := a.accountRepository.IsDeletedAccountByTelegramID(ctx, telegramInitModel.TelegramUser.ID)
	if isDeletedAccountWithTelegramID {
		log.Error("account not exist in db", logger.F("telegram_id", telegramInitModel.TelegramUser.ID))
		return nil, model.DuplicateAccountWithTelegramIDError
	}
	existAccountWithTelegramID, err := a.accountRepository.ExistsWithTelegramID(ctx, telegramInitModel.TelegramUser.ID)
	if err != nil {
		log.Error(
			"fail to check exist account with telegram id",
			logger.FError(err),
			logger.F("telegram_id", telegramInitModel.TelegramUser.ID),
		)
		return nil, err
	}
	if existAccountWithTelegramID {
		log.Error(
			"account already exist with the telegram id",
			logger.F("telegram_id", telegramInitModel.TelegramUser.ID),
		)
		return nil, model.DuplicateAccountWithTelegramIDError
	}
	gender, err := entity.GenderFromString(createAccount.Gender)
	if err != nil {
		log.Error("unknown gender from string", logger.FError(err))
		return nil, err
	}
	role, err := entity.RoleFromString(createAccount.Role)
	if err != nil {
		log.Error("unknown role from string", logger.FError(err))
		return nil, err
	}
	err = a.transactionProvider.Transact(func(composed transaction.ComposedRepository) error {
		var (
			companyID *int64
			err       error
		)
		if createAccount.HasCompany() {
			companyEntity := entity.Company{
				Name:        *createAccount.CompanyName,
				Description: *createAccount.CompanyDescription,
			}
			companyID, err = composed.Company.Create(ctx, &companyEntity)
			if err != nil {
				log.Error("fail to create company entity", logger.FError(err))
				return err
			}
			if companyID == nil {
				log.Error("company id has nil error")
				return model.NilError
			}
		}
		avatarFileHeader := createAccount.AvatarFileHeader
		documentFileHeader := createAccount.DocumentFileHeader
		if avatarFileHeader != nil {
			avatarAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(avatarFileHeader)
			if err != nil {
				log.Error("fail to upload and prepare a avatar attachment entity", logger.FError(err))
				return err
			}
			if avatarAttachmentEntity == nil {
				log.Error("avatar attachment has nil value", logger.FError(err))
				return model.NilError
			}
			avatarAttachmentEntityID, err := composed.Attachment.Create(ctx, avatarAttachmentEntity)
			if err != nil {
				log.Error("fail to record avatar attachment to db", logger.FError(err))
				return err
			}
			if avatarAttachmentEntityID == nil {
				log.Error("avatarAttachmentEntityID contains nil value")
				return model.NilError
			}
			avatarAttachmentEntity.ID = *avatarAttachmentEntityID
		}
		if documentFileHeader != nil {
			documentAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(documentFileHeader)
			if err != nil {
				log.Error("fail to upload and prepare a document attachment entity", logger.FError(err))
				return err
			}
			if documentAttachmentEntity == nil {
				log.Error("document attachment has nil value", logger.FError(err))
				return model.NilError
			}
			documentAttachmentEntityID, err := composed.Attachment.Create(ctx, documentAttachmentEntity)
			if err != nil {
				log.Error("fail to record document attachment to db", logger.FError(err))
				return err
			}
			if documentAttachmentEntityID == nil {
				log.Error("documentAttachmentEntityID contains nil value")
				return model.NilError
			}
			documentAttachmentEntity.ID = *documentAttachmentEntityID
		}
		var (
			avatarAttachmentEntityID   *int64
			documentAttachmentEntityID *int64
		)
		if avatarAttachmentEntity != nil {
			avatarAttachmentEntityID = &avatarAttachmentEntity.ID
		}
		if documentAttachmentEntity != nil {
			documentAttachmentEntityID = &documentAttachmentEntity.ID
		}
		accountEntity := entity.Account{
			TelegramID:           telegramInitModel.TelegramUser.ID,
			FirstName:            createAccount.FirstName,
			MiddleName:           createAccount.MiddleName,
			LastName:             createAccount.LastName,
			Role:                 role,
			Nickname:             createAccount.Nickname,
			AboutMe:              createAccount.AboutMe,
			Gender:               gender,
			Country:              &createAccount.Country,
			Location:             &createAccount.Location,
			CompanyID:            companyID,
			DocumentAttachmentID: documentAttachmentEntityID,
			AvatarAttachmentID:   avatarAttachmentEntityID,
		}
		accountID, err = composed.Account.Create(ctx, &accountEntity)
		if err != nil {
			log.Error("fail to record account in db", logger.FError(err))
			return err
		}
		if createAccount.HasTags() {
			tags := convertTags(*createAccount.Tags)
			for _, tag := range tags {
				tagEntity := entity.Tag{
					Title: tag,
				}
				_, err = composed.Tag.Create(ctx, &tagEntity, *accountID)
				if err != nil {
					log.Error("fail to create/add tag to account", logger.FError(err))
					return err
				}
			}
		}
		return nil
	})
	if err != nil && avatarAttachmentEntity != nil {
		_ = a.cleanupFileStore(avatarAttachmentEntity.FileName)
	}
	if err != nil && documentAttachmentEntity != nil {
		_ = a.cleanupFileStore(documentAttachmentEntity.FileName)
	}
	if err != nil {
		log.Error("fail to process creating a account", logger.FError(err))
		return nil, err
	}
	if accountID == nil {
		log.Error("account id has nil value", logger.FError(err))
		return nil, model.NilError
	}
	return accountID, nil
}

func (a *account) EditAccount(ctx context.Context, editAccount model.EditAccount) error {
	log := a.container.GetLogger()
	gender, err := entity.GenderFromString(editAccount.Gender)
	if err != nil {
		log.Error("unknown gender from string", logger.FError(err))
		return err
	}
	role, err := entity.RoleFromString(editAccount.Role)
	if err != nil {
		log.Error("unknown role from string", logger.FError(err))
		return err
	}
	var (
		newAvatarAttachmentEntity   *entity.Attachment
		removeAvatarFilename        *string
		newDocumentAttachmentEntity *entity.Attachment
		removeDocumentFilename      *string
	)
	err = a.transactionProvider.Transact(func(composed transaction.ComposedRepository) error {
		account, err := composed.Account.GetFullDetailByID(ctx, editAccount.ID)
		if err != nil {
			log.Error("fail to get account by id", logger.FError(err))
			switch err {
			case sql.ErrNoRows:
				return model.EntityNotFoundError
			default:
				return err
			}
		}
		var (
			oldAvatarFileName   *string
			oldDocumentFileName *string
		)
		if account.AvatarAttachment != nil {
			oldAvatarFileName = &account.AvatarAttachment.FileName
		}
		if account.DocumentAttachment != nil {
			oldDocumentFileName = &account.DocumentAttachment.FileName
		}
		if account.HasCompany() && editAccount.HasCompany() {
			company := account.Company
			if name := editAccount.CompanyName; name != nil {
				company.Name = *name
			}
			if description := editAccount.CompanyDescription; description != nil {
				company.Description = *description
			}
			if err := composed.Company.Update(ctx, company); err != nil {
				log.Error("fail to update company by id", logger.FError(err))
				return err
			}
		} else if account.HasCompany() && !editAccount.HasCompany() {
			if err := composed.Company.Delete(ctx, *account.CompanyID); err != nil {
				log.Error("fail to delete company by id", logger.FError(err))
				return err
			}
		} else if !account.HasCompany() && editAccount.HasCompany() {
			var company = entity.Company{}
			if name := editAccount.CompanyName; name != nil {
				company.Name = *name
			}
			if description := editAccount.CompanyDescription; description != nil {
				company.Description = *description
			}
			companyID, err := composed.Company.Create(ctx, &company)
			if err != nil {
				log.Error(
					"fail to create company for account",
					logger.F("account_id", editAccount.ID),
					logger.FError(err),
				)
				return err
			}
			newCompany, err := composed.Company.GetByID(ctx, *companyID)
			if err != nil {
				log.Error("fail to retrieve company", logger.F("company_id", *companyID), logger.FError(err))
				return err
			}
			account.CompanyID = companyID
			account.Company = newCompany
		}
		if editAccount.AvatarFileHeader != nil {
			newAvatarAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(editAccount.AvatarFileHeader)
			if err != nil {
				log.Error(
					"fail to upload avatar attachment and prepare entity for db",
					logger.FError(err),
				)
				return err
			}
			if account.AvatarAttachmentID != nil {
				newAvatarAttachmentEntity.ID = *account.AvatarAttachmentID
				if err := a.attachmentRepository.Update(ctx, newAvatarAttachmentEntity); err != nil {
					log.Error(
						"fail to update a avatar attachment in db",
						logger.F("avatar_id", newAvatarAttachmentEntity.ID),
						logger.FError(err),
					)
					return err
				}
			} else {
				avatarAttachmentID, err := a.attachmentRepository.Create(ctx, newAvatarAttachmentEntity)
				if err != nil {
					log.Error(
						"fail to create a avatar attachment in db",
						logger.FError(err),
					)
					return err
				}
				account.AvatarAttachmentID = avatarAttachmentID
			}
		}
		if editAccount.DocumentFileHeader != nil {
			newDocumentAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(editAccount.DocumentFileHeader)
			if err != nil {
				log.Error("fail to upload document attachment and prepare entity for db", logger.FError(err))
				return err
			}
			if account.DocumentAttachmentID != nil {
				newDocumentAttachmentEntity.ID = *account.DocumentAttachmentID
				if err := a.attachmentRepository.Update(ctx, newDocumentAttachmentEntity); err != nil {
					log.Error(
						"fail to update a document attachment in db",
						logger.F("document_id", newDocumentAttachmentEntity.ID),
						logger.FError(err),
					)
					return err
				}
			} else {
				documentAttachmentID, err := a.attachmentRepository.Create(ctx, newDocumentAttachmentEntity)
				if err != nil {
					log.Error(
						"fail to create a document attachment in db",
						logger.FError(err),
					)
					return err
				}
				account.DocumentAttachmentID = documentAttachmentID
			}
		}
		account.FirstName = editAccount.FirstName
		account.MiddleName = editAccount.MiddleName
		account.LastName = editAccount.LastName
		account.Role = role
		account.Nickname = editAccount.Nickname
		account.AboutMe = editAccount.AboutMe
		account.Gender = gender
		account.Country = &editAccount.Country
		account.Location = &editAccount.Location
		if err := composed.Account.Update(ctx, account); err != nil {
			log.Error("fail to update entity", logger.FError(err))
			return err
		}
		if oldAvatarFileName != nil {
			removeAvatarFilename = oldAvatarFileName
		}
		if oldDocumentFileName != nil {
			removeDocumentFilename = oldDocumentFileName
		}
		return nil
	})
	if err != nil && newAvatarAttachmentEntity != nil {
		_ = a.cleanupFileStore(newAvatarAttachmentEntity.FileName)
	}
	if err != nil && newDocumentAttachmentEntity != nil {
		_ = a.cleanupFileStore(newDocumentAttachmentEntity.FileName)
	}
	if err == nil && removeAvatarFilename != nil {
		_ = a.cleanupFileStore(*removeAvatarFilename)
	}
	if err == nil && removeDocumentFilename != nil {
		_ = a.cleanupFileStore(*removeDocumentFilename)
	}
	if err != nil {
		log.Error("fail to perform transaction while editing account", logger.FError(err))
		return err
	}
	log.Debug("we have successfully updated account |=)", logger.F("account_id", editAccount.ID))
	return nil
}

func (a *account) GeneratePairToken(accountID int64) (*model.PairToken, error) {
	log := a.container.GetLogger()
	pairToken, err := jwt.NewJWT(a.container.GetJWTSecretKey()).GeneratePairJWT(accountID)
	if err != nil {
		log.Error("fail to generate pair token", logger.FError(err))
		return nil, err
	}
	return &model.PairToken{
		Access:  pairToken.Access,
		Refresh: pairToken.Refresh,
	}, nil
}

func (a *account) ParseAccessToken(accessToken string) (*int64, error) {
	log := a.container.GetLogger()
	accessClaimsToken, err := jwt.NewJWT(a.container.GetJWTSecretKey()).ParseAccessToken(accessToken)
	if err != nil {
		log.Error("fail to parse/validate access token", logger.FError(err))
		return nil, err
	}
	if accessClaimsToken == nil {
		err := model.NilError
		log.Error("fail to get nil error", logger.FError(err))
		return nil, err
	}
	return &accessClaimsToken.ID, err
}

func (a *account) AuthenticationTelegram(ctx context.Context, telegramInitData string) (*int64, error) {
	log := a.container.GetLogger()
	var initData = telegram.InitData{
		Token: a.container.GetTelegramBotToken(),
	}
	telegramInitModel, err := initData.Decode(telegramInitData)
	if err != nil {
		log.Error("fail to decode telegram init data from string", logger.FError(err))
		return nil, err
	}
	ok, err := initData.Validate(telegramInitModel)
	if err != nil {
		log.Error("fail to process validate telegram init data", logger.FError(err))
		return nil, err
	}
	if !ok {
		err = model.InvalidTelegramInitDataError
		log.Error("invalid telegram initialization data provided", logger.FError(err))
		return nil, err
	}
	accountEntity, err := a.accountRepository.GetByTelegramID(ctx, telegramInitModel.TelegramUser.ID)
	if err != nil {
		log.Error("can't retrieve telegram by id", logger.FError(err))
		switch err {
		case sql.ErrNoRows:
			return nil, model.EntityNotFoundError
		default:
			return nil, err
		}
	}
	if accountEntity == nil {
		err = model.NilError
		log.Error("fail to get account by id", logger.FError(err))
		return nil, nil
	}
	return &accountEntity.ID, nil
}

func (a *account) GetDetailsAccount(ctx context.Context, id int64) (*model.Account, error) {
	log := a.container.GetLogger()
	account, err := a.accountRepository.GetFullDetailByID(ctx, id)
	if err != nil {
		log.Error("fail to get account by id", logger.FError(err))
		switch err {
		case sql.ErrNoRows:
			return nil, model.EntityNotFoundError
		default:
			return nil, err
		}
	}
	tags, err := a.tagRepository.GetTagsByAccountID(ctx, id)
	if err != nil {
		return nil, err
	}
	tagModels := make([]model.Tag, 0, len(tags))
	for _, tag := range tags {
		tagModel := converter.ConvertEntity2TagModel(&tag)
		tagModels = append(tagModels, *tagModel)
	}
	accountModel := converter.ConvertEntity2AccountModel(account)
	accountModel.Tags = &tagModels
	return accountModel, nil
}

func (a *account) AccountHasRole(ctx context.Context, accountID int64, role model.Role) (bool, error) {
	log := a.container.GetLogger()
	account, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		log.Error("fail to get by id", logger.FError(err))
		return false, nil
	}
	return account.Role.String() == string(role), nil
}

func (a *account) DeleteAccount(ctx context.Context, accountID int64) error {
	log := a.container.GetLogger()
	account, err := a.accountRepository.GetFullDetailByID(ctx, accountID)
	if err != nil {
		log.Error("fail to get profile by id", logger.FError(err))
		switch err {
		case sql.ErrNoRows:
			return model.EntityNotFoundError
		default:
			return err
		}
	}
	err = a.transactionProvider.Transact(func(composed transaction.ComposedRepository) error {
		if err = composed.Account.Delete(ctx, account.ID); err != nil {
			log.Error("fail to delete account from db", logger.FError(err))
			return err
		}
		if companyID := account.CompanyID; companyID != nil {
			err = composed.Company.Delete(ctx, *companyID)
			if err != nil {
				log.Error("fail to delete company from db", logger.FError(err))
				return err
			}
		}
		if avatarAttachmentID := account.AvatarAttachmentID; avatarAttachmentID != nil {
			err = composed.Attachment.Delete(ctx, *avatarAttachmentID)
			if err != nil {
				log.Error("fail to delete avatar attachment from db", logger.FError(err))
				return err
			}
		}
		if documentAttachmentID := account.DocumentAttachmentID; documentAttachmentID != nil {
			err := composed.Attachment.Delete(ctx, *documentAttachmentID)
			if err != nil {
				log.Error("fail to delete document attachment from db", logger.FError(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error("fail to delete a account and associate data to it from db", logger.FError(err))
		return err
	}
	return nil
}
func (a *account) ChangeRole(ctx context.Context, accountID int64, role model.Role) error {
	log := a.container.GetLogger()
	accountEntity, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.FError(err))
		return err
	}
	accountEntity.Role, err = entity.RoleFromString(string(role))
	if err != nil {
		log.Error("fail to parse role", logger.FError(err))
		return err
	}
	if err = a.accountRepository.Update(ctx, accountEntity); err != nil {
		log.Error("fail to update account", logger.FError(err))
		return err
	}
	return nil
}

func (a *account) uploadAndPrepareAttachmentEntity(fileHeader *multipart.FileHeader) (*entity.Attachment, error) {
	log := a.container.GetLogger()
	attachmentUploadFile, err := a.unpackFileHeader(fileHeader)
	if err != nil {
		log.Error("fail to unpack a file header", logger.FError(err))
		return nil, err
	}
	if attachmentUploadFile == nil {
		log.Error("avatar upload file has nil error")
		return nil, model.NilError
	}
	attachmentFilePath, err := a.fileStorage.UploadFile(attachmentUploadFile.Name, attachmentUploadFile.File)
	if err != nil {
		log.Error("fail to upload avatar file to file storage", logger.FError(err))
		return nil, err
	}
	attachment := entity.Attachment{
		FileName: attachmentUploadFile.Name,
		Path:     attachmentFilePath,
	}
	return &attachment, nil
}

func (a *account) cleanupFileStore(name string) error {
	return a.fileStorage.DeleteFile(name)
}

func (a *account) unpackFileHeader(fileHeader *multipart.FileHeader) (*uploadFile, error) {
	log := a.container.GetLogger()
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("fail to open avatar file", logger.FError(err))
		return nil, err
	}
	fileExt, err := utils.ExtFromFileName(fileHeader.Filename)
	if err != nil {
		log.Error("fail to extract extension from filename", logger.FError(err))
		return nil, err
	}
	fileName := uuid.NewString()
	return &uploadFile{
		Name: fmt.Sprintf("%s%s", fileName, *fileExt),
		File: file,
	}, nil
}

func convertTags(tags []string) []string {
	convertedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		editedTag := strings.TrimSpace(strings.ToLower(tag))
		convertedTags = append(convertedTags, editedTag)
	}
	return convertedTags
}
