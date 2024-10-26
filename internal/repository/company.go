package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain"
	"time"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) (*int64, error)
	FetchByID(ctx context.Context, id int64) (*domain.Company, error)
}

type companyRepository struct {
	conn *sql.DB
}

func NewCompanyRepository(conn *sql.DB) CompanyRepository {
	return &companyRepository{
		conn: conn,
	}
}

func (c *companyRepository) Create(ctx context.Context, company *domain.Company) (*int64, error) {
	var id int64
	err := c.conn.QueryRowContext(ctx,
		"INSERT INTO company (name, description, created_at) VALUES ($1, $2, $3) RETURNING id",
		company.Name,
		company.Description,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func (c *companyRepository) FetchByID(ctx context.Context, id int64) (*domain.Company, error) {
	row := c.conn.QueryRowContext(ctx, "SELECT name, description FROM company WHERE id = $1", id)
	var company domain.Company
	err := row.Scan(company.Name, company.Description)
	company.ID = &id
	return &company, err
}
