package usecase

import (
	"context"
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
)

type Account interface {
	CreateAccount(ctx context.Context, createAccount model.CreateAccount) (*int64, error)
	GeneratePairToken(accountID int64) (*model.PairToken, error)
	AuthenticationTelegram(ctx context.Context, telegramInitData string) (*int64, error)
	ParseAccessToken(accessToken string) (*int64, error)
	GetDetailsAccount(ctx context.Context, id int64) (*model.Account, error)
	EditAccount(ctx context.Context, editAccount model.EditAccount) error
	GetMatchAccounts(ctx context.Context, accountID int64, limit int64) ([]model.Account, error)
}

type account struct {
	container            container.Container
	fileStorage          filestorage.FileStorage
	accountRepository    repository.Account
	attachmentRepository repository.Attachment
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
	transactionProvider *transaction.Provider,
) Account {
	return &account{
		container:            container,
		fileStorage:          fileStorage,
		accountRepository:    accountRepository,
		attachmentRepository: attachmentRepository,
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
		avatarAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(createAccount.AvatarFileHeader)
		if err != nil {
			log.Error("fail to upload and prepare a avatar attachment entity", logger.FError(err))
			return err
		}
		if avatarAttachmentEntity == nil {
			log.Error("avatar attachment has nil value", logger.FError(err))
			return model.NilError
		}
		documentAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(createAccount.DocumentFileHeader)
		if err != nil {
			log.Error("fail to upload and prepare a document attachment entity", logger.FError(err))
			return err
		}
		if documentAttachmentEntity == nil {
			log.Error("document attachment has nil value", logger.FError(err))
			return model.NilError
		}
		avatarAttachmentEntityID, err := composed.Attachment.Create(ctx, avatarAttachmentEntity)
		if err != nil {
			log.Error("fail to record avatar attachment to db", logger.FError(err))
			return err
		}
		documentAttachmentEntityID, err := composed.Attachment.Create(ctx, documentAttachmentEntity)
		if err != nil {
			log.Error("fail to record document attachment to db", logger.FError(err))
			return err
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
			DocumentAttachmentID: avatarAttachmentEntityID,
			AvatarAttachmentID:   documentAttachmentEntityID,
		}
		accountID, err = composed.Account.Create(ctx, &accountEntity)
		if err != nil {
			log.Error("fail to record account in db", logger.FError(err))
			return err
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
			return err
		}
		if account.AvatarAttachmentID == nil || account.DocumentAttachmentID == nil {
			log.Error("both avatar and document are required for edit a account")
			return model.NilError
		}
		var oldAvatarFileName = account.AvatarAttachment.FileName
		var oldDocumentFileName = account.DocumentAttachment.FileName
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
		newAvatarAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(editAccount.AvatarFileHeader)
		if err != nil {
			log.Error(
				"fail to upload avatar attachment and prepare entity for db",
				logger.FError(err),
			)
			return err
		}
		newAvatarAttachmentEntity.ID = *account.AvatarAttachmentID
		if err := a.attachmentRepository.Update(ctx, newAvatarAttachmentEntity); err != nil {
			log.Error(
				"fail to update a avatar attachment in db",
				logger.F("avatar_id", newAvatarAttachmentEntity.ID),
				logger.FError(err),
			)
			return err
		}
		newDocumentAttachmentEntity, err = a.uploadAndPrepareAttachmentEntity(editAccount.DocumentFileHeader)
		if err != nil {
			log.Error("fail to upload document attachment and prepare entity for db", logger.FError(err))
			return err
		}
		newDocumentAttachmentEntity.ID = *account.DocumentAttachmentID
		if err := a.attachmentRepository.Update(ctx, newDocumentAttachmentEntity); err != nil {
			log.Error(
				"fail to update a document attachment in db",
				logger.F("document_id", newDocumentAttachmentEntity.ID),
				logger.FError(err),
			)
			return err
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
		removeAvatarFilename = &oldAvatarFileName
		removeDocumentFilename = &oldDocumentFileName
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
		return nil, err
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
		log.Error("fail to get account by id")
		return nil, err
	}
	accountModel := converter.ConvertEntity2AccountModel(account)
	return accountModel, nil
}

func (a *account) GetMatchAccounts(ctx context.Context, accountID int64, limit int64) ([]model.Account, error) {
	log := a.container.GetLogger()
	accountEntity, err := a.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.FError(err))
		return nil, err
	}
	accountEntities, err := a.accountRepository.GetMatchedAccounts(ctx, accountID, accountEntity.Role.Opposite(), limit)
	if err != nil {
		log.Error("fail to get matched accounts", logger.FError(err))
		return nil, err
	}
	accounts := make([]model.Account, 0, len(accountEntities))
	for _, accountEntity := range accountEntities {
		err := a.accountRepository.SeenAccount(ctx, accountID, accountEntity.ID)
		if err != nil {
			log.Error("fail to mark account as seen", logger.FError(err))
			return nil, err
		}
		account := converter.ConvertEntity2AccountModel(&accountEntity)
		accounts = append(accounts, *account)
	}
	return accounts, nil
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
