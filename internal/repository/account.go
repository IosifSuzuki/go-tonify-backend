package repository

import (
	"database/sql"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/model"
	"golang.org/x/net/context"
	"time"
)

type AccountRepository interface {
	ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error)
	Create(ctx context.Context, account *domain.Account) (*int64, error)
	FetchByID(ctx context.Context, id int64) (*domain.Account, error)
	FetchByTelegramID(ctx context.Context, telegramID int64) (*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	GetMatchedAccounts(ctx context.Context, accountID int64, role model.Role, limit int64) ([]domain.Account, error)
	ExistsSeenAccount(ctx context.Context, viewerAccountID, viewedAccountID int64) (bool, error)
	SeenAccount(ctx context.Context, viewerAccountID, viewedAccountID int64) error
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

func (a *accountRepository) GetMatchedAccounts(ctx context.Context, accountID int64, role model.Role, limit int64) ([]domain.Account, error) {
	query := "SELECT" +
		"	account.id, " +
		"	account.telegram_id, " +
		"	account.first_name, " +
		"	account.middle_name, " +
		"	account.last_name, " +
		"	account.nickname, " +
		"	account.role, " +
		"	account.about_me, " +
		"	account.gender, " +
		"	account.country, " +
		"	account.location, " +
		"	account.avatar_id, " +
		"	account.document_id, " +
		"	account.company_id, " +
		"	account.created_at, " +
		"	account.updated_at " +
		"FROM" +
		"	account " +
		"LEFT JOIN account_seen ON account.id = account_seen.viewed_account_id" +
		"	AND account_seen.viewer_account_id = $1 " +
		"WHERE" +
		"	account.role = $2" +
		"	AND account.id != $3 " +
		"ORDER BY" +
		"	account_seen.rating ASC " +
		"LIMIT $4;"
	rows, err := a.conn.QueryContext(ctx, query, accountID, role, accountID, limit)
	if err != nil {
		return nil, err
	}
	accountDomains := make([]domain.Account, 0, limit)
	for rows.Next() {
		var middleName sql.NullString
		var nickname sql.NullString
		var aboutMe sql.NullString
		var companyID sql.NullInt64
		var country sql.NullString
		var location sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		var avatarID sql.NullInt64
		var documentID sql.NullInt64
		var accountDomain domain.Account
		err = rows.Scan(
			&accountDomain.ID,
			&accountDomain.TelegramID,
			&accountDomain.FirstName,
			&middleName,
			&accountDomain.LastName,
			&nickname,
			&accountDomain.Role,
			&aboutMe,
			&accountDomain.Gender,
			&country,
			&location,
			&avatarID,
			&documentID,
			&companyID,
			&createdAt,
			&updatedAt,
		)
		if middleName.Valid {
			accountDomain.MiddleName = &middleName.String
		}
		if aboutMe.Valid {
			accountDomain.AboutMe = &aboutMe.String
		}
		if nickname.Valid {
			accountDomain.Nickname = &nickname.String
		}
		if companyID.Valid {
			accountDomain.CompanyID = &companyID.Int64
		}
		if createdAt.Valid {
			accountDomain.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			accountDomain.UpdatedAt = &updatedAt.Time
		}
		if country.Valid {
			accountDomain.Country = &country.String
		}
		if location.Valid {
			accountDomain.Location = &location.String
		}
		if avatarID.Valid {
			accountDomain.AvatarAttachmentID = &avatarID.Int64
		}
		if documentID.Valid {
			accountDomain.DocumentAttachmentID = &documentID.Int64
		}
		if err != nil {
			return nil, err
		}
		accountDomains = append(accountDomains, accountDomain)
	}
	return accountDomains, nil
}

func (a *accountRepository) ExistsSeenAccount(ctx context.Context, viewerAccountID int64, viewedAccountID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM account_seen WHERE viewer_account_id = $1 AND viewed_account_id = $2);"
	var exists bool
	err := a.conn.QueryRowContext(ctx, query, viewerAccountID, viewedAccountID).Scan(&exists)
	return exists, err
}

func (a *accountRepository) SeenAccount(ctx context.Context, viewerAccountID int64, viewedAccountID int64) error {
	exist, err := a.ExistsSeenAccount(ctx, viewerAccountID, viewedAccountID)
	if err != nil {
		return err
	}
	if exist {
		query := "UPDATE account_seen SET rating = rating + 1 " +
			"WHERE viewer_account_id = $1 AND viewed_account_id = $2"
		_, err := a.conn.ExecContext(ctx, query, viewerAccountID, viewedAccountID)
		if err != nil {
			return err
		}
		return nil
	}

	query := "INSERT INTO account_seen (viewer_account_id, viewed_account_id, rating) VALUES ($1, $2, $3);"
	_, err = a.conn.ExecContext(ctx, query, viewerAccountID, viewedAccountID, 0)
	if err != nil {
		return err
	}
	return nil
}
