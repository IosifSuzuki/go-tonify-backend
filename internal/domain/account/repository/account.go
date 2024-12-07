package repository

import (
	"database/sql"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
	"golang.org/x/net/context"
	"time"
)

type Account interface {
	ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error)
	Create(ctx context.Context, account *entity.Account) (*int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	GetByIDWithCompany(ctx context.Context, id int64) (*entity.Account, error)
	GetFullDetailByID(ctx context.Context, id int64) (*entity.Account, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id int64) error
	GetMatchedAccounts(ctx context.Context, accountID int64, role entity.Role, limit int64) ([]entity.Account, error)
	ExistsSeenAccount(ctx context.Context, viewerAccountID, viewedAccountID int64) (bool, error)
	SeenAccount(ctx context.Context, viewerAccountID, viewedAccountID int64) error
}

type account struct {
	conn psql.Operation
}

func NewAccount(conn psql.Operation) Account {
	return &account{
		conn: conn,
	}
}

func (a *account) ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM account WHERE telegram_id=$1 AND account.deleted_at IS NULL);"
	var exists bool
	err := a.conn.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	return exists, err
}

func (a *account) Create(ctx context.Context, account *entity.Account) (*int64, error) {
	query := "INSERT INTO account (" +
		"	telegram_id, " +
		"	first_name, " +
		"	middle_name, " +
		"	last_name, " +
		"	nickname, " +
		"	role, " +
		"	about_me, " +
		"	gender, " +
		"	country, " +
		"	location, " +
		"	company_id, " +
		"	avatar_id, " +
		"	document_id, " +
		"	created_at" +
		") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) " +
		"RETURNING id;"
	var id int64
	err := a.conn.QueryRowContext(
		ctx,
		query,
		account.TelegramID,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.Role.String(),
		account.AboutMe,
		account.Gender.String(),
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

func (a *account) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	query := "SELECT " +
		"	telegram_id, " +
		"	first_name, " +
		"	middle_name, " +
		"	last_name, " +
		"	nickname, " +
		"	role, " +
		"	about_me, " +
		"	gender, " +
		"	country, " +
		"	location, " +
		"	avatar_id, " +
		"	document_id, " +
		"	company_id, " +
		"	created_at," +
		"	updated_at " +
		"FROM account WHERE id = $1 AND deleted_at IS NULL;"
	row := a.conn.QueryRowContext(ctx, query, id)
	var (
		middleName sql.NullString
		aboutMe    sql.NullString
		nickname   sql.NullString
		companyID  sql.NullInt64
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
		avatarID   sql.NullInt64
		documentID sql.NullInt64
		gender     string
		role       string
	)
	var account = entity.Account{
		ID:       id,
		Country:  new(string),
		Location: new(string),
	}
	err := row.Scan(
		&account.TelegramID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&role,
		&aboutMe,
		&gender,
		account.Country,
		account.Location,
		&avatarID,
		&documentID,
		&companyID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	account.Role, _ = entity.RoleFromString(role)
	account.Gender, _ = entity.GenderFromString(gender)
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
	return &account, err
}

func (a *account) GetByIDWithCompany(ctx context.Context, id int64) (*entity.Account, error) {
	query := "SELECT" +
		"	telegram_id, " +
		"	first_name," +
		"	middle_name," +
		"	last_name," +
		"	nickname," +
		"	role," +
		"	about_me," +
		"	gender," +
		"	country," +
		"	location," +
		"	avatar_id," +
		"	document_id," +
		"	company_id," +
		"	company.name," +
		"	company.description," +
		"	account.created_at," +
		"	account.updated_at " +
		"FROM account " +
		"	JOIN company ON account.company_id = company.id " +
		"WHERE account.id = $1 AND account.deleted_at IS NULL;"
	var account = entity.Account{
		ID:       id,
		Country:  new(string),
		Location: new(string),
	}
	row := a.conn.QueryRowContext(ctx, query, id)
	var (
		middleName         sql.NullString
		aboutMe            sql.NullString
		nickname           sql.NullString
		avatarID           sql.NullInt64
		documentID         sql.NullInt64
		companyID          sql.NullInt64
		companyName        sql.NullString
		companyDescription sql.NullString
		createdAt          sql.NullTime
		updatedAt          sql.NullTime
		gender             string
		role               string
	)
	err := row.Scan(
		&account.TelegramID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&role,
		&aboutMe,
		&gender,
		&account.Country,
		&account.Location,
		&avatarID,
		&documentID,
		&companyID,
		&companyName,
		&companyDescription,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	account.Role, _ = entity.RoleFromString(role)
	account.Gender, _ = entity.GenderFromString(gender)
	if middleName.Valid {
		account.MiddleName = &middleName.String
	}
	if nickname.Valid {
		account.Nickname = &nickname.String
	}
	if aboutMe.Valid {
		account.AboutMe = &aboutMe.String
	}
	if avatarID.Valid {
		account.AvatarAttachmentID = &avatarID.Int64
	}
	if documentID.Valid {
		account.DocumentAttachmentID = &documentID.Int64
	}
	if companyID.Valid {
		var company = new(entity.Company)
		company.ID = companyID.Int64
		account.Company = company
		account.CompanyID = &companyID.Int64
	}
	if companyName.Valid {
		account.Company.Name = companyName.String
	}
	if companyDescription.Valid {
		account.Company.Description = companyDescription.String
	}
	if createdAt.Valid {
		account.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		account.UpdatedAt = &updatedAt.Time
	}
	return &account, nil
}

func (a *account) GetFullDetailByID(ctx context.Context, id int64) (*entity.Account, error) {
	query := "SELECT" +
		"	telegram_id, " +
		"	first_name," +
		"	middle_name," +
		"	last_name," +
		"	nickname," +
		"	role," +
		"	about_me," +
		"	gender," +
		"	country," +
		"	location," +
		"	avatar_id," +
		"	avatar.file_name," +
		"	avatar.path," +
		"	avatar.created_at," +
		"	avatar.updated_at," +
		"	document_id," +
		"	document.file_name," +
		"	document.path," +
		"	document.created_at," +
		"	document.updated_at," +
		"	company_id," +
		"	company.name," +
		"	company.description," +
		"	company.created_at," +
		"	company.updated_at," +
		"	account.created_at," +
		"	account.updated_at " +
		"FROM account " +
		"	JOIN company ON account.company_id = company.id " +
		"	JOIN attachment as avatar ON account.avatar_id = avatar.id " +
		"	JOIN attachment as document ON account.document_id = document.id " +
		"WHERE account.id = $1 AND account.deleted_at IS NULL;"
	row := a.conn.QueryRowContext(ctx, query, id)
	var (
		middleName         sql.NullString
		aboutMe            sql.NullString
		nickname           sql.NullString
		companyID          sql.NullInt64
		companyName        sql.NullString
		companyDescription sql.NullString
		companyCreatedAt   sql.NullTime
		companyUpdatedAt   sql.NullTime
		createdAt          sql.NullTime
		updatedAt          sql.NullTime
		avatarID           sql.NullInt64
		avatarFileName     sql.NullString
		avatarPath         sql.NullString
		avatarCreatedAt    sql.NullTime
		avatarUpdatedAt    sql.NullTime
		documentID         sql.NullInt64
		documentFileName   sql.NullString
		documentPath       sql.NullString
		documentCreatedAt  sql.NullTime
		documentUpdatedAt  sql.NullTime
		role               string
		gender             string
	)
	var account = entity.Account{
		ID:       id,
		Country:  new(string),
		Location: new(string),
	}
	err := row.Scan(
		&account.TelegramID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&role,
		&aboutMe,
		&gender,
		&account.Country,
		&account.Location,
		&avatarID,
		&avatarFileName,
		&avatarPath,
		&avatarCreatedAt,
		&avatarUpdatedAt,
		&documentID,
		&documentFileName,
		&documentPath,
		&documentCreatedAt,
		&documentUpdatedAt,
		&companyID,
		&companyName,
		&companyDescription,
		&companyCreatedAt,
		&companyUpdatedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	account.Role, _ = entity.RoleFromString(role)
	account.Gender, _ = entity.GenderFromString(gender)
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
		var company = new(entity.Company)
		company.ID = companyID.Int64
		account.CompanyID = &companyID.Int64
		account.Company = company
	}
	if companyName.Valid {
		account.Company.Name = companyName.String
	}
	if companyDescription.Valid {
		account.Company.Description = companyDescription.String
	}
	if companyCreatedAt.Valid {
		account.Company.CreatedAt = &companyCreatedAt.Time
	}
	if companyUpdatedAt.Valid {
		account.Company.UpdatedAt = &companyUpdatedAt.Time
	}
	if createdAt.Valid {
		account.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		account.UpdatedAt = &updatedAt.Time
	}
	if avatarID.Valid {
		var avatarAttachment = new(entity.Attachment)
		avatarAttachment.ID = avatarID.Int64
		account.AvatarAttachmentID = &avatarID.Int64
		account.AvatarAttachment = avatarAttachment
	}
	if avatarFileName.Valid && account.AvatarAttachment != nil {
		account.AvatarAttachment.FileName = avatarFileName.String
	}
	if avatarPath.Valid && account.AvatarAttachment != nil {
		account.AvatarAttachment.Path = &avatarPath.String
	}
	if avatarCreatedAt.Valid && account.AvatarAttachment != nil {
		account.AvatarAttachment.CreatedAt = &avatarCreatedAt.Time
	}
	if avatarUpdatedAt.Valid && account.AvatarAttachment != nil {
		account.AvatarAttachment.UpdatedAt = &avatarUpdatedAt.Time
	}
	if documentID.Valid {
		var documentAttachment = new(entity.Attachment)
		documentAttachment.ID = documentID.Int64
		account.DocumentAttachmentID = &documentID.Int64
		account.DocumentAttachment = documentAttachment
	}
	if documentFileName.Valid && account.DocumentAttachment != nil {
		account.DocumentAttachment.FileName = documentFileName.String
	}
	if documentPath.Valid && account.DocumentAttachment != nil {
		account.DocumentAttachment.Path = &documentPath.String
	}
	if documentCreatedAt.Valid && account.DocumentAttachment != nil {
		account.DocumentAttachment.CreatedAt = &documentCreatedAt.Time
	}
	if documentUpdatedAt.Valid && account.DocumentAttachment != nil {
		account.DocumentAttachment.UpdatedAt = &documentUpdatedAt.Time
	}
	return &account, nil
}

func (a *account) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.Account, error) {
	query := "SELECT " +
		"	id, " +
		"	first_name, " +
		"	middle_name, " +
		"	last_name, " +
		"	nickname, " +
		"	role, " +
		"	about_me, " +
		"	gender, " +
		"	country, " +
		"	location, " +
		"	avatar_id, " +
		"	document_id, " +
		"	company_id, " +
		"	created_at, " +
		"	updated_at " +
		"FROM account WHERE telegram_id = $1 AND deleted_at IS NULL;"
	row := a.conn.QueryRowContext(ctx, query, telegramID)
	var (
		middleName sql.NullString
		nickname   sql.NullString
		aboutMe    sql.NullString
		companyID  sql.NullInt64
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
		avatarID   sql.NullInt64
		documentID sql.NullInt64
		gender     string
		role       string
	)
	var account = entity.Account{
		TelegramID: telegramID,
	}
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&middleName,
		&account.LastName,
		&nickname,
		&role,
		&aboutMe,
		&gender,
		account.Country,
		account.Location,
		account.AvatarAttachmentID,
		account.DocumentAttachmentID,
		&companyID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	account.Role, _ = entity.RoleFromString(role)
	account.Gender, _ = entity.GenderFromString(gender)
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
	return &account, err
}

func (a *account) Update(ctx context.Context, account *entity.Account) error {
	query := "UPDATE account SET " +
		"	first_name = $1, " +
		"	middle_name = $2, " +
		"	last_name = $3, " +
		"	nickname = $4, " +
		"	role = $5, " +
		"	about_me = $6, " +
		"	gender = $7, " +
		"	country = $8, " +
		"	location = $9, " +
		"	avatar_id = $10, " +
		"	document_id = $11, " +
		"	company_id = $12, " +
		"	updated_at = $13 " +
		"WHERE id = $14;"
	_, err := a.conn.ExecContext(
		ctx,
		query,
		account.FirstName,
		account.MiddleName,
		account.LastName,
		account.Nickname,
		account.Role.String(),
		account.AboutMe,
		account.Gender.String(),
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

func (a *account) Delete(ctx context.Context, id int64) error {
	query := "UPDATE account SET" +
		"	deleted_at = $1" +
		"WHERE id = $2;"
	_, err := a.conn.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (a *account) GetMatchedAccounts(ctx context.Context, accountID int64, role entity.Role, limit int64) ([]entity.Account, error) {
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
		"	avatar.file_name," +
		"	avatar.path," +
		"	avatar.created_at," +
		"	avatar.updated_at," +
		"	account.document_id," +
		"	account.company_id, " +
		"	account.created_at, " +
		"	account.updated_at " +
		"FROM" +
		"	account " +
		"LEFT JOIN account_seen ON account.id = account_seen.viewed_account_id" +
		"	AND account_seen.viewer_account_id = $1 " +
		"LEFT JOIN attachment as avatar ON account.avatar_id = avatar.id " +
		"WHERE" +
		"	account.role = $2" +
		"	AND account.id != $3 " +
		"ORDER BY" +
		"	account_seen.rating ASC " +
		"LIMIT $4;"
	rows, err := a.conn.QueryContext(ctx, query, accountID, role.String(), accountID, limit)
	if err != nil {
		return nil, err
	}
	accounts := make([]entity.Account, 0, limit)
	for rows.Next() {
		var (
			middleName      sql.NullString
			nickname        sql.NullString
			aboutMe         sql.NullString
			companyID       sql.NullInt64
			country         sql.NullString
			location        sql.NullString
			createdAt       sql.NullTime
			updatedAt       sql.NullTime
			avatarID        sql.NullInt64
			avatarName      sql.NullString
			avatarPath      sql.NullString
			avatarCreatedAt sql.NullTime
			avatarUpdatedAt sql.NullTime
			documentID      sql.NullInt64
			role            string
			gender          string
		)
		var account entity.Account
		err = rows.Scan(
			&account.ID,
			&account.TelegramID,
			&account.FirstName,
			&middleName,
			&account.LastName,
			&nickname,
			&role,
			&aboutMe,
			&gender,
			&country,
			&location,
			&avatarID,
			&avatarName,
			&avatarPath,
			&avatarCreatedAt,
			&avatarUpdatedAt,
			&documentID,
			&companyID,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		account.Role, _ = entity.RoleFromString(role)
		account.Gender, _ = entity.GenderFromString(gender)
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
		if country.Valid {
			account.Country = &country.String
		}
		if location.Valid {
			account.Location = &location.String
		}
		if avatarID.Valid {
			var attachment = entity.Attachment{
				ID:        avatarID.Int64,
				Path:      new(string),
				CreatedAt: new(time.Time),
				UpdatedAt: new(time.Time),
			}
			account.AvatarAttachment = &attachment
		}
		if avatarName.Valid {
			account.AvatarAttachment.FileName = avatarName.String
		}
		if avatarPath.Valid {
			account.AvatarAttachment.Path = &avatarPath.String
		}
		if avatarCreatedAt.Valid {
			account.AvatarAttachment.CreatedAt = &avatarCreatedAt.Time
		}
		if avatarUpdatedAt.Valid {
			account.AvatarAttachment.UpdatedAt = &avatarUpdatedAt.Time
		}
		if documentID.Valid {
			account.DocumentAttachmentID = &documentID.Int64
		}
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (a *account) ExistsSeenAccount(ctx context.Context, viewerAccountID int64, viewedAccountID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM account_seen WHERE viewer_account_id = $1 AND viewed_account_id = $2);"
	var exists bool
	err := a.conn.QueryRowContext(ctx, query, viewerAccountID, viewedAccountID).Scan(&exists)
	return exists, err
}

func (a *account) SeenAccount(ctx context.Context, viewerAccountID int64, viewedAccountID int64) error {
	exist, err := a.ExistsSeenAccount(ctx, viewerAccountID, viewedAccountID)
	if err != nil {
		return err
	}
	if exist {
		query := "UPDATE account_seen SET " +
			"	rating = rating + 1 " +
			"WHERE viewer_account_id = $1 AND viewed_account_id = $2;"
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
