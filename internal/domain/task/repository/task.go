package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
)

type Task interface {
	Create(ctx context.Context, task *entity.Task) (*int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Task, error)
	CountByID(ctx context.Context, id int64) (*int64, error)
	GetList(ctx context.Context, ownerID int64, offset int64, limit int64) ([]entity.Task, error)
}

type task struct {
	conn psql.Operation
}

func NewTask(conn psql.Operation) Task {
	return &task{
		conn: conn,
	}
}

func (t *task) Create(ctx context.Context, task *entity.Task) (*int64, error) {
	var id int64
	query := "INSERT INTO task (" +
		"	owner_id, " +
		"	title, " +
		"	description " +
		") VALUES ($1, $2, $3) " +
		"RETURNING id;"
	err := t.conn.QueryRowContext(ctx, query, task.OwnerID, task.Title, task.Description).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (t *task) GetByID(ctx context.Context, id int64) (*entity.Task, error) {
	query := "SELECT " +
		"	owner_id, " +
		"	title, " +
		"	description, " +
		"	created_at, " +
		"	updated_at " +
		"FROM task " +
		"	WHERE id = $1 AND deleted_at IS NULL;"
	var (
		task      entity.Task
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)
	task.ID = id
	err := t.conn.QueryRowContext(ctx, query, id).Scan(
		&task.OwnerID,
		&task.Title,
		&task.Description,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		task.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		task.UpdatedAt = &updatedAt.Time
	}
	return &task, nil
}

func (t *task) CountByID(ctx context.Context, id int64) (*int64, error) {
	query := "SELECT COUNT(*) " +
		"FROM task " +
		"	WHERE owner_id = $1;"
	var count int64
	err := t.conn.QueryRowContext(ctx, query, id).Scan(
		&count,
	)
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (t *task) GetList(ctx context.Context, ownerID int64, offset int64, limit int64) ([]entity.Task, error) {
	query := "SELECT " +
		"	id, " +
		"	title, " +
		"	description, " +
		"	created_at, " +
		"	updated_at " +
		"FROM task " +
		"	WHERE owner_id = $1 AND deleted_at IS NULL " +
		"LIMIT $2 " +
		"OFFSET $3;"
	rows, err := t.conn.QueryContext(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, err
	}
	tasks := make([]entity.Task, 0, limit)
	for rows.Next() {
		var (
			createdAt sql.NullTime
			updatedAt sql.NullTime
		)
		var task entity.Task
		task.OwnerID = ownerID
		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			task.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			task.UpdatedAt = &updatedAt.Time
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
