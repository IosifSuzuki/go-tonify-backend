package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
	"time"
)

type Attachment interface {
	Create(ctx context.Context, attachment *entity.Attachment) (*int64, error)
	Update(ctx context.Context, attachment *entity.Attachment) error
	GetByID(ctx context.Context, id int64) (*entity.Attachment, error)
	Delete(ctx context.Context, id int64) error
}

type attachment struct {
	conn psql.Operation
}

func NewAttachment(conn psql.Operation) Attachment {
	return &attachment{
		conn: conn,
	}
}

func (a *attachment) Create(ctx context.Context, attachment *entity.Attachment) (*int64, error) {
	query := "INSERT INTO attachment (" +
		"	file_name, " +
		"	path, " +
		"	created_at" +
		") VALUES ($1, $2, $3) RETURNING id;"
	var id int64
	err := a.conn.QueryRowContext(
		ctx,
		query,
		attachment.FileName,
		attachment.Path,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (a *attachment) Update(ctx context.Context, attachment *entity.Attachment) error {
	query := "UPDATE attachment SET " +
		"	file_name = $1, " +
		"	path = $2, " +
		"	updated_at = $3 " +
		"	WHERE id = $4;"
	_, err := a.conn.ExecContext(
		ctx,
		query,
		attachment.FileName,
		attachment.Path,
		time.Now(),
		attachment.ID,
	)
	return err
}

func (a *attachment) GetByID(ctx context.Context, id int64) (*entity.Attachment, error) {
	query := "SELECT " +
		"	file_name, " +
		"	path, " +
		"	created_at, " +
		"	updated_at " +
		"FROM attachment WHERE id = $1 AND delete_at IS NOT NULL;"
	var (
		path      sql.NullString
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)
	attachment := entity.Attachment{
		ID:        id,
		Path:      new(string),
		CreatedAt: new(time.Time),
		UpdatedAt: new(time.Time),
	}
	row := a.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&attachment.FileName,
		&path,
		&createdAt,
		&updatedAt,
	)
	if path.Valid {
		attachment.Path = &path.String
	}
	if createdAt.Valid {
		attachment.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		attachment.UpdatedAt = &updatedAt.Time
	}
	return &attachment, err
}

func (a *attachment) Delete(ctx context.Context, id int64) error {
	query := "UPDATE attachment SET " +
		"	deleted_at = $1 " +
		"WHERE id = $2;"
	_, err := a.conn.ExecContext(ctx,
		query,
		time.Now(),
		id,
	)
	return err
}
