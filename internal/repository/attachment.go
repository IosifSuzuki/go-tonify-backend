package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain"
	"time"
)

type attachmentRepository struct {
	conn *sql.DB
}

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *domain.Attachment) (*int64, error)
	Update(ctx context.Context, attachment *domain.Attachment) error
	FetchByID(ctx context.Context, id int64) (*domain.Attachment, error)
	Delete(ctx context.Context, id int64) error
}

func NewAttachment(conn *sql.DB) AttachmentRepository {
	return &attachmentRepository{
		conn: conn,
	}
}

func (a *attachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) (*int64, error) {
	query := "INSERT INTO attachment (file_name, status, created_at) VALUES ($1, $2, $3) RETURNING id;"
	var id int64
	err := a.conn.QueryRowContext(
		ctx,
		query,
		attachment.FileName,
		attachment.Status,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (a *attachmentRepository) Update(ctx context.Context, attachment *domain.Attachment) error {
	query := "UPDATE attachment SET file_name = $1, path = $2, status = $3, updated_at = $4 WHERE id = $5"
	_, err := a.conn.ExecContext(ctx, query, attachment.FileName, attachment.Path, attachment.Status, attachment.UpdatedAt, attachment.ID)
	return err
}

func (a *attachmentRepository) FetchByID(ctx context.Context, id int64) (*domain.Attachment, error) {
	query := "SELECT file_name, path, status, created_at, updated_at FROM attachment WHERE id = $1"
	var path sql.NullString
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	attachment := domain.NewAttachment()
	row := a.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&attachment.FileName,
		&path,
		&attachment.Status,
		&createdAt,
		&updatedAt,
	)
	attachment.ID = id
	if path.Valid {
		attachment.Path = &path.String
	}
	if createdAt.Valid {
		attachment.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		attachment.UpdatedAt = &updatedAt.Time
	}
	return attachment, err
}

func (a *attachmentRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM attachment WHERE id = $1"
	_, err := a.conn.ExecContext(ctx, query, id)
	return err
}
