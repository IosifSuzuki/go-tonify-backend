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
	query := "SELECT telegram_id, first_name, middle_name, last_name, nickname, about_me, gender, country, location, company_id, created_at, updated_at, deleted_at FROM account WHERE id = $1"
	row := a.conn.QueryRowContext(ctx, query, id)
	var account domain.Account
	account.ID = &id
	err := row.Scan(
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
		account.CreatedAt,
		account.UpdatedAt,
		account.DeletedAt,
	)
	return &account, err
}
