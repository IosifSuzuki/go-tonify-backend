package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/utils/encrypt"
	"go-tonify-backend/pkg/logger"
	"net/url"
	"strconv"
)

type AuthService interface {
	CreateAccount(ctx context.Context, createAccount *model.CreateAccount) (*int64, error)
	AuthorizationAccount(ctx context.Context, credential *model.Credential) (*model.Account, error)
	GenerateAccountJWT(ctx context.Context, accountID int64) (*model.PairToken, error)
	ParseAccessAccountJWT(token string) (*model.AccessClaimsToken, error)
	ParseRefreshAccountJWT(token string) (*model.RefreshClaimsToken, error)
}

type authService struct {
	container         container.Container
	accountRepository repository.AccountRepository
	companyRepository repository.CompanyRepository
	jwt               *encrypt.JWT
}

func NewAuthService(
	container container.Container,
	accountRepository repository.AccountRepository,
	companyRepository repository.CompanyRepository,
) AuthService {
	return &authService{
		container:         container,
		accountRepository: accountRepository,
		companyRepository: companyRepository,
		jwt:               encrypt.NewJWT(container),
	}
}

func (a *authService) CreateAccount(ctx context.Context, createAccount *model.CreateAccount) (*int64, error) {
	log := a.container.GetLogger()
	ctx, cancel := context.WithTimeout(ctx, a.container.GetContentTimeout())
	defer cancel()
	telegramInitData, err := DecodeTelegramInitData(createAccount.TelegramRawInitData)
	if err != nil {
		log.Error("fail to decode telegram init fata", logger.FError(err))
		return nil, model.TelegramInitDataDecodeError
	}
	isValidTelegramInitData, err := a.validateTelegramInitData(
		telegramInitData,
		a.container.GetTelegramConfig().BotToken,
	)
	if !isValidTelegramInitData {
		log.Error("invalid init telegram data", logger.FError(err))
		err = model.TelegramInitDataValidationError
		return nil, err
	} else if err != nil {
		log.Error("fail to processing validating init telegram data", logger.FError(err))
		return nil, model.TelegramInitDataValidationError
	}
	accountExists, err := a.accountRepository.ExistsWithTelegramID(ctx, telegramInitData.TelegramUser.ID)
	if err != nil {
		log.Error("failed to check existence of Telegram ID in system", logger.FError(err))
		return nil, model.DataBaseOperationError
	} else if accountExists {
		log.Error("account already exist in system", logger.FError(err))
		return nil, model.AccountAlreadyExistsError
	}
	gender := string(createAccount.Gender)
	accountEntity := domain.Account{
		TelegramID:           telegramInitData.TelegramUser.ID,
		FirstName:            createAccount.FirstName,
		MiddleName:           createAccount.MiddleName,
		LastName:             createAccount.LastName,
		Role:                 string(createAccount.Role),
		Nickname:             createAccount.Nickname,
		AboutMe:              createAccount.AboutMe,
		Gender:               gender,
		Country:              &createAccount.Country,
		Location:             &createAccount.Location,
		CompanyID:            createAccount.CompanyID,
		AvatarAttachmentID:   createAccount.AvatarID,
		DocumentAttachmentID: createAccount.DocumentID,
	}
	accountID, err := a.accountRepository.Create(ctx, &accountEntity)
	if err != nil {
		return nil, model.DataBaseOperationError
	}
	return accountID, nil
}

func (a *authService) GenerateAccountJWT(ctx context.Context, accountID int64) (*model.PairToken, error) {
	accessClaimsToken := model.AccessClaimsToken{
		ID: accountID,
	}
	refreshClaimsToken := model.RefreshClaimsToken{
		ID: accountID,
	}
	accessToken, err := a.jwt.GenerateToken(accessClaimsToken)
	if err != nil {
		return nil, model.CreationJWTError
	}
	refreshToken, err := a.jwt.GenerateToken(refreshClaimsToken)
	if err != nil {
		return nil, model.CreationJWTError
	}
	return &model.PairToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (a *authService) ParseAccessAccountJWT(token string) (*model.AccessClaimsToken, error) {
	return a.jwt.ParseAccessToken(token)
}

func (a *authService) ParseRefreshAccountJWT(token string) (*model.RefreshClaimsToken, error) {
	return a.jwt.ParseRefreshToken(token)
}

func (a *authService) AuthorizationAccount(ctx context.Context, credential *model.Credential) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, a.container.GetContentTimeout())
	defer cancel()
	telegramInitData, err := DecodeTelegramInitData(credential.TelegramRawInitData)
	if err != nil {
		return nil, model.TelegramInitDataDecodeError
	}
	isValidTelegramInitData, err := a.validateTelegramInitData(
		telegramInitData,
		a.container.GetTelegramConfig().BotToken,
	)
	if !isValidTelegramInitData {
		return nil, model.TelegramInitDataValidationError
	} else if err != nil {
		return nil, model.TelegramInitDataValidationError
	}
	accountExists, err := a.accountRepository.ExistsWithTelegramID(ctx, telegramInitData.TelegramUser.ID)
	if !accountExists {
		return nil, model.AccountNotExistsError
	} else if err != nil {
		return nil, model.DataBaseOperationError
	}
	account, err := a.accountRepository.FetchByTelegramID(ctx, telegramInitData.TelegramUser.ID)
	if err != nil {
		return nil, model.DataBaseOperationError
	}
	return &model.Account{
		ID:         account.ID,
		TelegramID: account.TelegramID,
		FirstName:  account.FirstName,
		MiddleName: account.MiddleName,
		LastName:   account.LastName,
		Nickname:   account.Nickname,
		Role:       model.Role(account.Role),
		AboutMe:    account.AboutMe,
		Gender:     model.NewGender(account.Gender),
		Country:    account.Country,
		Location:   account.Location,
		CompanyID:  account.CompanyID,
	}, nil
}

func DecodeTelegramInitData(data string) (*model.TelegramInitData, error) {
	values, err := url.ParseQuery(data)
	if err != nil {
		return nil, err
	}
	queryIDs := values["query_id"]
	if len(queryIDs) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	queryID := queryIDs[0]
	users := values["user"]
	if len(users) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	payloadUser := users[0]
	var telegramUser model.TelegramUser
	if err := json.Unmarshal([]byte(payloadUser), &telegramUser); err != nil {
		return nil, model.TelegramInitDataDecodeError
	}
	authDates := values["auth_date"]
	if len(authDates) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	authDate, err := strconv.Atoi(authDates[0])
	if err != nil {
		return nil, model.TelegramInitDataDecodeError
	}
	hashes := values["hash"]
	if len(hashes) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	hash := hashes[0]
	return &model.TelegramInitData{
		QueryID:             queryID,
		TelegramUserPayload: payloadUser,
		TelegramUser:        telegramUser,
		AuthDate:            uint(authDate),
		Hash:                hash,
	}, nil
}

func (a *authService) validateTelegramInitData(telegramInitData *model.TelegramInitData, token string) (bool, error) {
	log := a.container.GetLogger()
	dataCheckString := fmt.Sprintf(
		"auth_date=%d\nquery_id=%s\nuser=%s",
		telegramInitData.AuthDate,
		telegramInitData.QueryID,
		telegramInitData.TelegramUserPayload,
	)
	log.Debug("start validate", logger.F("dataCheckString", dataCheckString))
	secretKey, err := encrypt.GetSHA256Signature([]byte(token), []byte("WebAppData"))
	if err != nil {
		log.Error("fail to secret key", logger.FError(err))
		return false, err
	}
	generatedHash, err := encrypt.GetSHA256Signature([]byte(dataCheckString), secretKey)
	if err != nil {
		log.Error("fail to generate hash for validation data", logger.FError(err))
		return false, err
	}
	generatedHexHash := hex.EncodeToString(generatedHash)
	if generatedHexHash == telegramInitData.Hash {
		return true, nil
	}
	log.Error(
		"corrupted data",
		logger.F("secret_key", secretKey),
		logger.F("hash from init data", telegramInitData.Hash),
		logger.F("generated hash for check equation", generatedHexHash),
	)
	return true, nil
}
