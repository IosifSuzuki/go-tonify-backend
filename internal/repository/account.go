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
	query := "INSERT INTO account (telegram_id, first_name, middle_name, last_name, nickname, role, about_me, gender, country, location, company_id, avatar_id, document_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id;"
	var id int64
	err := a.conn.QueryRowContext(
		ctx,
		query,
		account.TelegramID,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.Role,
		account.AboutMe,
		account.Gender,
		account.Country,
		account.Location,
		account.CompanyID,
		account.AvatarAttachmentID,
		account.DocumentAttachmentID,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func (a *accountRepository) FetchByID(ctx context.Context, id int64) (*domain.Account, error) {
	query := "SELECT telegram_id, first_name, middle_name, last_name, nickname, role, about_me, gender, country, location, avatar_id, document_id, company_id, created_at FROM account WHERE id = $1"
	row := a.conn.QueryRowContext(ctx, query, id)
	var middleName sql.NullString
	var aboutMe sql.NullString
	var nickname sql.NullString
	var companyID sql.NullInt64
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var avatarID sql.NullInt64
	var documentID sql.NullInt64
	account := domain.NewAccount()
	err := row.Scan(
		&account.TelegramID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&account.Role,
		&aboutMe,
		&account.Gender,
		account.Country,
		account.Location,
		&avatarID,
		&documentID,
		&companyID,
		&createdAt,
	)
	if middleName.Valid {
		account.MiddleName = &middleName.String
	}
	if aboutMe.Valid {
		account.AboutMe = &aboutMe.String
	}
	if nickname.Valid {
		account.Nickname = &nickname.String
	}
	if companyID.Valid {
		account.CompanyID = &companyID.Int64
	}
	if createdAt.Valid {
		account.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		account.UpdatedAt = &updatedAt.Time
	}
	if avatarID.Valid {
		account.AvatarAttachmentID = &avatarID.Int64
	}
	if documentID.Valid {
		account.DocumentAttachmentID = &documentID.Int64
	}
	account.ID = id
	return account, err
}

func (a *accountRepository) FetchByTelegramID(ctx context.Context, telegramID int64) (*domain.Account, error) {
	query := "SELECT id, first_name, middle_name, last_name, nickname, role, about_me, gender, country, location, avatar_id, document_id, company_id, created_at, updated_at FROM account WHERE telegram_id = $1"
	row := a.conn.QueryRowContext(ctx, query, telegramID)
	var middleName sql.NullString
	var nickname sql.NullString
	var aboutMe sql.NullString
	var companyID sql.NullInt64
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var avatarID sql.NullInt64
	var documentID sql.NullInt64
	account := domain.NewAccount()
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&account.Role,
		&aboutMe,
		&account.Gender,
		account.Country,
		account.Location,
		account.AvatarAttachmentID,
		account.DocumentAttachmentID,
		&companyID,
		&createdAt,
		&updatedAt,
	)
	if middleName.Valid {
		account.MiddleName = &middleName.String
	}
	if aboutMe.Valid {
		account.AboutMe = &aboutMe.String
	}
	if nickname.Valid {
		account.Nickname = &nickname.String
	}
	if companyID.Valid {
		account.CompanyID = &companyID.Int64
	}
	if createdAt.Valid {
		account.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		account.UpdatedAt = &updatedAt.Time
	}
	if avatarID.Valid {
		account.AvatarAttachmentID = &avatarID.Int64
	}
	if documentID.Valid {
		account.DocumentAttachmentID = &documentID.Int64
	}
	account.TelegramID = telegramID
	return account, err
}

func (a *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	query := "UPDATE account SET first_name = $1, middle_name = $2, last_name = $3, nickname = $4, role = $5, about_me = $6, gender = $7, " +
		"country = $8, location = $9, avatar_id = $10, document_id = $11, company_id = $12, updated_at = $13 WHERE id = $14"
	_, err := a.conn.ExecContext(
		ctx,
		query,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.Role,
		account.AboutMe,
		account.Gender,
		account.Country,
		account.Location,
		account.AvatarAttachmentID,
		account.DocumentAttachmentID,
		account.CompanyID,
		time.Now(),
		account.ID,
	)
	return err
}
