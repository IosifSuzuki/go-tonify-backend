package repository

import (
	"context"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
)

type Category interface {
	GetCategoriesByAccountID(ctx context.Context, accountID int64) ([]entity.Category, error)
	GetCategoriesByTaskID(ctx context.Context, taskID int64) ([]entity.Category, error)
	AddCategoryToAccount(ctx context.Context, categoryID int64, accountID int64) error
	DeleteCategoriesFromAccount(ctx context.Context, accountID int64) error
	DeleteCategoriesFromTask(ctx context.Context, taskID int64) error
	AddCategoryToTask(ctx context.Context, categoryID int64, taskID int64) error
	GetAll(ctx context.Context, offset int64, limit int64) ([]entity.Category, error)
	GetAllNumberRows(ctx context.Context) (*int64, error)
}

type category struct {
	conn psql.Operation
}

func NewCategory(conn psql.Operation) Category {
	return &category{
		conn: conn,
	}
}

func (c *category) GetCategoriesByAccountID(ctx context.Context, accountID int64) ([]entity.Category, error) {
	query := "SELECT category.id, category.title FROM category " +
		"	JOIN account_category " +
		"	ON account_category.category_id = category.id " +
		"	WHERE account_category.account_id = $1;"
	rows, err := c.conn.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	categories := make([]entity.Category, 0, 0)
	for rows.Next() {
		var category entity.Category
		err = rows.Scan(
			&category.ID,
			&category.Title,
		)
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *category) GetCategoriesByTaskID(ctx context.Context, taskID int64) ([]entity.Category, error) {
	query := "SELECT category.id, category.title FROM category " +
		"	JOIN account_task " +
		"	ON account_task.task_id = task.id " +
		"	WHERE account_task.task_id = $1;"
	rows, err := c.conn.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	categories := make([]entity.Category, 0, 0)
	for rows.Next() {
		var category entity.Category
		err = rows.Scan(
			&category.ID,
			&category.Title,
		)
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *category) DeleteCategoriesFromAccount(ctx context.Context, accountID int64) error {
	query := "DELETE FROM account_category WHERE account_id = $1"
	_, err := c.conn.ExecContext(ctx, query, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (c *category) DeleteCategoriesFromTask(ctx context.Context, taskID int64) error {
	query := "DELETE FROM task_category WHERE task_id = $1"
	_, err := c.conn.ExecContext(ctx, query, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (c *category) AddCategoryToAccount(ctx context.Context, categoryID int64, accountID int64) error {
	query := "INSERT INTO account_category (category_id, account_id) VALUES ($1, $2);"
	_, err := c.conn.ExecContext(ctx, query, categoryID, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (c *category) AddCategoryToTask(ctx context.Context, categoryID int64, taskID int64) error {
	query := "INSERT INTO task_category (category_id, task_id) VALUES ($1, $2);"
	_, err := c.conn.ExecContext(ctx, query, categoryID, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (c *category) GetAll(ctx context.Context, offset int64, limit int64) ([]entity.Category, error) {
	query := "SELECT id, title FROM category LIMIT $1 OFFSET $2;"
	rows, err := c.conn.QueryContext(ctx, query, limit, offset)
	categories := make([]entity.Category, 0, 0)
	for rows.Next() {
		var category entity.Category
		err = rows.Scan(
			&category.ID,
			&category.Title,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *category) GetAllNumberRows(ctx context.Context) (*int64, error) {
	query := "SELECT COUNT(*) as all_rows FROM category;"
	var count int64
	if err := c.conn.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return nil, err
	}
	return &count, nil
}
