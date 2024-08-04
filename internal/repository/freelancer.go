package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain"
	"time"
)

type FreelancerRepository interface {
	Create(ctx context.Context, company *domain.Freelancer) (*int64, error)
	ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error)
}

type freelancerRepository struct {
	conn *sql.DB
}

func NewFreelancerRepository(conn *sql.DB) FreelancerRepository {
	return &freelancerRepository{
		conn: conn,
	}
}

func (f *freelancerRepository) ExistsWithTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM freelancer WHERE telegram_id=$1);"
	var exists bool
	err := f.conn.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	return exists, err
}

func (f *freelancerRepository) Create(ctx context.Context, freelancer *domain.Freelancer) (*int64, error) {
	query := "INSERT INTO freelancer (telegram_id, first_name, middle_name, last_name, gender, country, city, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	var id int64
	err := f.conn.QueryRowContext(
		ctx,
		query,
		freelancer.TelegramID,
		freelancer.FirstName,
		freelancer.MiddleName,
		freelancer.LastName,
		freelancer.Gender,
		freelancer.Country,
		freelancer.City,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, err
}
