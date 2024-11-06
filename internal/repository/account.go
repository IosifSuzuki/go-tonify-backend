package repository

import (
	"database/sql"
	"go-tonify-backend/internal/domain"
	"golang.org/x/net/context"
	"time"
)

type AccountRepository interface {
	ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error)
	Create(ctx context.Context, account *domain.Account) (*int64, error)
	FetchByID(ctx context.Context, id int64) (*domain.Account, error)
	FetchByTelegramID(ctx context.Context, telegramID int64) (*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
}

type accountRepository struct {
	conn *sql.DB
}

func NewAccountRepository(conn *sql.DB) AccountRepository {
	return &accountRepository{
		conn: conn,
	}
}

func (a *accountRepository) ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM account WHERE telegram_id=$1);"
	var exists bool
	err := a.conn.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	return exists, err
}

func (a *accountRepository) Create(ctx context.Context, account *domain.Account) (*int64, error) {
	query := "INSERT INTO account (telegram_id, first_name, middle_name, last_name, nickname, about_me, gender, country, location, company_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;"
	var id int64
	err := a.conn.QueryRowContext(
		ctx,
		query,
		account.TelegramID,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.AboutMe,
		account.Gender,
		account.Country,
		account.Location,
		account.CompanyID,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func (a *accountRepository) FetchByID(ctx context.Context, id int64) (*domain.Account, error) {
	query := "SELECT telegram_id, first_name, middle_name, last_name, nickname, about_me, gender, country, location, company_id FROM account WHERE id = $1"
	row := a.conn.QueryRowContext(ctx, query, id)
	var middleName sql.NullString
	var aboutMe sql.NullString
	var nickname sql.NullString
	var companyID sql.NullInt64
	account := domain.NewAccount()
	err := row.Scan(
		account.TelegramID,
		account.FirstName,
		&middleName,
		account.LastName,
		&nickname,
		&aboutMe,
		account.Gender,
		account.Country,
		account.Location,
		&companyID,
	)
	account.ID = &id
	return account, err
}

func (a *accountRepository) FetchByTelegramID(ctx context.Context, telegramID int64) (*domain.Account, error) {
	query := "SELECT id, first_name, middle_name, last_name, nickname, about_me, gender, country, location, company_id, created_at FROM account WHERE telegram_id = $1"
	row := a.conn.QueryRowContext(ctx, query, telegramID)
	var middleName sql.NullString
	var nickname sql.NullString
	var aboutMe sql.NullString
	var companyID sql.NullInt64
	account := domain.NewAccount()
	err := row.Scan(
		account.ID,
		account.FirstName,
		&middleName,
		account.LastName,
		&nickname,
		&aboutMe,
		account.Gender,
		account.Country,
		account.Location,
		&companyID,
		account.CreatedAt,
	)
	account.TelegramID = &telegramID
	return account, err
}

func (a *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	query := "UPDATE account SET first_name = $1, middle_name = $2, last_name = $3, nickname = $4, about_me = $5, gender = $6, " +
		"country = $7, location = $8, updated_at = $9 WHERE id = $10"
	_, err := a.conn.ExecContext(
		ctx,
		query,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.AboutMe,
		account.Gender,
		account.Country,
		account.Location,
		account.UpdatedAt,
		account.ID,
	)
	return err
}
