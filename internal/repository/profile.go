package repository

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/domain"
)

type ProfileRepository interface {
	GetByID(c context.Context, id string) (*domain.Profile, error)
}

type profileRepository struct {
	conn *sql.DB
}

func NewProfileRepository(conn *sql.DB) ProfileRepository {
	return &profileRepository{
		conn: conn,
	}
}

func (p *profileRepository) GetByID(_ context.Context, _ string) (*domain.Profile, error) {
	return &domain.Profile{
		ID:        "1234-5678",
		FirstName: "Andriy",
		LastName:  "Melnyk",
		UserName:  "@test",
	}, nil
}
