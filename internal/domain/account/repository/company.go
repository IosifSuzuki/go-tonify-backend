package repository

import (
	"context"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/pkg/psql"
	"time"
)

type Company interface {
	Create(ctx context.Context, company *entity.Company) (*int64, error)
	Update(ctx context.Context, company *entity.Company) error
	GetByID(ctx context.Context, id int64) (*entity.Company, error)
	Delete(ctx context.Context, id int64) error
}

type company struct {
	conn psql.Operation
}

func NewCompany(conn psql.Operation) Company {
	return &company{
		conn: conn,
	}
}

func (c *company) Create(ctx context.Context, company *entity.Company) (*int64, error) {
	var id int64
	err := c.conn.QueryRowContext(
		ctx,
		"INSERT INTO company ("+
			"	name, "+
			"	description, "+
			"	created_at "+
			") VALUES ($1, $2, $3) "+
			"RETURNING id;",
		company.Name,
		company.Description,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func (c *company) Update(ctx context.Context, company *entity.Company) error {
	query := "UPDATE company SET " +
		"	name = $1, " +
		"	description = $2, " +
		"	updated_at = $3 " +
		"WHERE id = $4;"
	_, err := c.conn.ExecContext(
		ctx,
		query,
		company.Name,
		company.Description,
		time.Now(),
		company.ID,
	)
	return err
}

func (c *company) GetByID(ctx context.Context, id int64) (*entity.Company, error) {
	row := c.conn.QueryRowContext(ctx, "SELECT "+
		"	name, "+
		"	description, "+
		"	created_at, "+
		"	updated_at "+
		"FROM company WHERE id = $1 AND deleted_at IS NOT NULL;",
		id,
	)
	company := entity.Company{
		ID:        id,
		CreatedAt: new(time.Time),
		UpdatedAt: new(time.Time),
	}
	err := row.Scan(
		&company.Name,
		&company.Description,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	return &company, err
}

func (c *company) Delete(ctx context.Context, id int64) error {
	query := "UPDATE company SET " +
		"	deleted_at = $1 " +
		"WHERE id = $2;"
	_, err := c.conn.ExecContext(ctx, query, time.Now(), id)
	return err
}
