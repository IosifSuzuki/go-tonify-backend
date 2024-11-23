package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain"
	"time"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) (*int64, error)
	Update(ctx context.Context, company *domain.Company) error
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

func (c *companyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := "UPDATE company SET name = $1, description = $2, updated_at = $3 WHERE id = $4"
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
func (c *companyRepository) FetchByID(ctx context.Context, id int64) (*domain.Company, error) {
	row := c.conn.QueryRowContext(ctx, "SELECT name, description FROM company WHERE id = $1", id)
	company := domain.NewCompany()
	err := row.Scan(company.Name, company.Description)
	company.ID = id
	return company, err
}
