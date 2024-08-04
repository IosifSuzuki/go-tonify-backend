package repository

import (
	"database/sql"
	"go-tonify-backend/internal/domain"
	"golang.org/x/net/context"
	"time"
)

type ClientRepository interface {
	ExistsWithID(ctx context.Context, telegramID int64) (bool, error)
	Create(ctx context.Context, client *domain.Client) (*int64, error)
	FetchByID(ctx context.Context, id int64) (*domain.Client, error)
}

type clientRepository struct {
	conn *sql.DB
}

func NewClientRepository(conn *sql.DB) ClientRepository {
	return &clientRepository{
		conn: conn,
	}
}

func (c *clientRepository) ExistsWithID(ctx context.Context, telegramID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM client WHERE telegram_id=$1);"
	var exists bool
	err := c.conn.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	return exists, err
}

func (c *clientRepository) Create(ctx context.Context, client *domain.Client) (*int64, error) {
	query := "INSERT INTO client (telegram_id, first_name, middle_name, last_name, gender, country, city, company_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;"
	var id int64
	err := c.conn.QueryRowContext(
		ctx,
		query,
		client.TelegramID,
		client.FirstName,
		client.MiddleName,
		client.LastName,
		client.Gender,
		client.Country,
		client.City,
		client.CompanyID,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func (c *clientRepository) FetchByID(ctx context.Context, id int64) (*domain.Client, error) {
	query := "SELECT telegram_id, first_name, middle_name, last_name, gender, country, city, company_id, created_at, updated_at, deleted_at FROM client WHERE id = $1"
	row := c.conn.QueryRowContext(ctx, query, id)
	var client domain.Client
	client.ID = &id
	err := row.Scan(
		client.TelegramID,
		client.FirstName,
		client.MiddleName,
		client.LastName,
		client.Gender,
		client.Country,
		client.City,
		client.CompanyID,
	)
	return &client, err
}
