package repository

import (
	"context"
	"go-tonify-backend/internal/domain"
)

type ProfileRepository interface {
	GetByID(c context.Context, id string) (*domain.Profile, error)
}

type profileRepository struct {
}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}

func (p *profileRepository) GetByID(_ context.Context, _ string) (*domain.Profile, error) {
	return &domain.Profile{
		ID:        "1234-5678",
		FirstName: "Andriy",
		LastName:  "Melnyk",
		UserName:  "@test",
	}, nil
}
