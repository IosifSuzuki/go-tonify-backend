package repository

import (
	"context"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
)

type Tag interface {
	Create(ctx context.Context, tag *entity.Tag, accountID int64) (*int64, error)
	ExistTagWithTitle(ctx context.Context, title string) (bool, error)
	GetTagByTitle(ctx context.Context, title string) (*entity.Tag, error)
	GetTagsByAccountID(ctx context.Context, accountID int64) ([]entity.Tag, error)
}

type tag struct {
	conn psql.Operation
}

func NewTag(conn psql.Operation) Tag {
	return &tag{
		conn: conn,
	}
}

func (t *tag) Create(ctx context.Context, tag *entity.Tag, accountID int64) (*int64, error) {
	exists, err := t.ExistTagWithTitle(ctx, tag.Title)
	if err != nil {
		return nil, err
	}
	var tagID int64
	if !exists {
		query := "INSERT INTO tag (title) VALUES ($1) RETURNING id;"
		if err := t.conn.QueryRowContext(ctx, query, tag.Title).Scan(&tagID); err != nil {
			return nil, err
		}
	} else {
		tag, err := t.GetTagByTitle(ctx, tag.Title)
		if err != nil {
			return nil, err
		}
		tagID = tag.ID
	}
	accountTag := entity.AccountTag{
		AccountID: accountID,
		TagID:     tagID,
	}
	query := "INSERT INTO account_tag (account_id, tag_id) VALUES ($1, $2)"
	_, err = t.conn.ExecContext(ctx, query, accountTag.AccountID, accountTag.TagID)
	if err != nil {
		return nil, err
	}
	return &tagID, nil
}

func (t *tag) GetTagsByAccountID(ctx context.Context, accountID int64) ([]entity.Tag, error) {
	query := "SELECT tag.id, tag.title " +
		"	FROM tag" +
		"	LEFT JOIN account_tag ON account_tag.tag_id = tag.id" +
		"	WHERE account_tag.account_id = $1;"
	rows, err := t.conn.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	tags := make([]entity.Tag, 0)
	for rows.Next() {
		var tag entity.Tag
		err = rows.Scan(
			&tag.ID,
			&tag.Title,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (t *tag) ExistTagWithTitle(ctx context.Context, title string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tag WHERE title=$1);"
	var exists bool
	if err := t.conn.QueryRowContext(ctx, query, title).Scan(&exists); err != nil {
		return false, nil
	}
	return exists, nil
}

func (t *tag) GetTagByTitle(ctx context.Context, title string) (*entity.Tag, error) {
	query := "SELECT id, title FROM tag WHERE title = $1;"
	var tag entity.Tag
	err := t.conn.QueryRowContext(ctx, query, title).Scan(
		&tag.ID,
		&tag.Title,
	)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
