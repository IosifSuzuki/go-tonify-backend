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
	"log"
	"net/url"
	"strconv"
)

type AuthService interface {
	CreateAccount(ctx context.Context, createAccount *model.CreateAccount) (*int64, error)
	GenerateAccountJWT(ctx context.Context, accountID int64) (*model.PairToken, error)
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
	ctx, cancel := context.WithTimeout(ctx, a.container.GetContentTimeout())
	defer cancel()
	telegramInitData, err := decodeTelegramInitData(createAccount.TelegramRawInitData)
	if err != nil {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	isValidTelegramInitData, err := validateTelegramInitData(
		telegramInitData,
		a.container.GetTelegramConfig().BotToken,
	)
	if !isValidTelegramInitData {
		return nil, model.TelegramInitDataValidationError
	} else if err != nil {
		log.Println(err)
		return nil, model.TelegramInitDataValidationError
	}
	accountExists, err := a.accountRepository.ExistsWithTelegramID(ctx, telegramInitData.TelegramUser.ID)
	if err != nil {
		log.Println(err)
		return nil, model.DataBaseOperationError
	} else if accountExists {
		return nil, model.AccountAlreadyExistsError
	}
	var companyID *int64
	if createAccount.CompanyName != nil && createAccount.CompanyDescription != nil {
		companyEntity := domain.Company{
			Name:        createAccount.CompanyName,
			Description: createAccount.CompanyDescription,
		}
		companyID, err = a.companyRepository.Create(ctx, &companyEntity)
		if err != nil {
			log.Println(err)
			return nil, model.DataBaseOperationError
		}
	}
	gender := string(createAccount.Gender)
	accountEntity := domain.Account{
		TelegramID: &telegramInitData.TelegramUser.ID,
		FirstName:  &createAccount.FirstName,
		MiddleName: createAccount.MiddleName,
		LastName:   &createAccount.LastName,
		Nickname:   createAccount.Nickname,
		AboutMe:    createAccount.AboutMe,
		Gender:     &gender,
		Country:    &createAccount.Country,
		Location:   &createAccount.Location,
		CompanyID:  companyID,
	}
	accountID, err := a.accountRepository.Create(ctx, &accountEntity)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, model.CreationJWTError
	}
	refreshToken, err := a.jwt.GenerateToken(refreshClaimsToken)
	if err != nil {
		log.Println(err)
		return nil, model.CreationJWTError
	}
	return &model.PairToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func decodeTelegramInitData(data string) (*model.TelegramInitData, error) {
	values, err := url.ParseQuery(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	queryIDs := values["query_id"]
	if len(queryIDs) != 1 {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	queryID := queryIDs[0]
	users := values["user"]
	if len(users) != 1 {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	payloadUser := users[0]
	var telegramUser model.TelegramUser
	if err := json.Unmarshal([]byte(payloadUser), &telegramUser); err != nil {
		log.Println("unmarshal", err)
		return nil, model.TelegramInitDataDecodeError
	}
	authDates := values["auth_date"]
	if len(authDates) != 1 {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	authDate, err := strconv.Atoi(authDates[0])
	if err != nil {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	hashes := values["hash"]
	if len(hashes) != 1 {
		log.Println(err)
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

func validateTelegramInitData(telegramInitData *model.TelegramInitData, token string) (bool, error) {
	dataCheckString := fmt.Sprintf(
		"auth_date=%d\nquery_id=%s\nuser=%s",
		telegramInitData.AuthDate,
		telegramInitData.QueryID,
		telegramInitData.TelegramUserPayload,
	)
	secretKey, err := encrypt.GetSHA256Signature([]byte(token), []byte("WebAppData"))
	if err != nil {
		log.Println(err)
		return false, err
	}
	generatedHash, err := encrypt.GetSHA256Signature([]byte(dataCheckString), secretKey)
	generatedHexHash := hex.EncodeToString(generatedHash)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return generatedHexHash == telegramInitData.Hash, nil
}
