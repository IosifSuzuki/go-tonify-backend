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
	CreateClient(ctx context.Context, createClient *model.CreateClient) (*int64, error)
	GenerateClientJWT(ctx context.Context, clientID int64) (*model.PairToken, error)
}

type authService struct {
	container         container.Container
	clientRepository  repository.ClientRepository
	companyRepository repository.CompanyRepository
	jwt               *encrypt.JWT
}

func NewAuthService(
	clientRepository repository.ClientRepository,
	companyRepository repository.CompanyRepository,
	container container.Container,
) AuthService {
	return &authService{
		container:         container,
		clientRepository:  clientRepository,
		companyRepository: companyRepository,
		jwt:               encrypt.NewJWT(container),
	}
}

func (a *authService) CreateClient(ctx context.Context, createClient *model.CreateClient) (*int64, error) {
	ctx, cancel := context.WithTimeout(ctx, a.container.GetContentTimeout())
	defer cancel()
	telegramInitData, err := decodeTelegramInitData(createClient.TelegramRawInitData)
	if err != nil {
		log.Println(err)
		return nil, model.TelegramInitDataDecodeError
	}
	clientExists, err := a.clientRepository.ExistsWithID(ctx, telegramInitData.TelegramUser.ID)
	if err != nil {
		log.Println(err)
		return nil, model.DataBaseOperationError
	} else if clientExists {
		return nil, model.ClientAlreadyExistsError
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
	companyEntity := domain.Company{
		Name:        &createClient.CompanyName,
		Description: &createClient.CompanyDescription,
	}
	companyID, err := a.companyRepository.Create(ctx, &companyEntity)
	if err != nil {
		log.Println(err)
		return nil, model.DataBaseOperationError
	}
	clientEntity := domain.Client{
		TelegramID: &telegramInitData.TelegramUser.ID,
		FirstName:  &createClient.FirstName,
		MiddleName: createClient.MiddleName,
		LastName:   &createClient.LastName,
		Gender:     &createClient.Gender,
		Country:    &createClient.Country,
		City:       &createClient.City,
		CompanyID:  companyID,
	}
	clientID, err := a.clientRepository.Create(ctx, &clientEntity)
	if err != nil {
		log.Println(err)
		return nil, model.DataBaseOperationError
	}
	return clientID, nil
}

func (a *authService) GenerateClientJWT(ctx context.Context, clientID int64) (*model.PairToken, error) {
	accessClaimsToken := model.AccessClaimsToken{
		ID:   clientID,
		Role: "client",
	}
	refreshClaimsToken := model.RefreshClaimsToken{
		ID: 0,
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
