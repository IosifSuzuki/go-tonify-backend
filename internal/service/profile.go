package service

import (
	"context"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/repository"
)

type ProfileService interface {
	GetProfileByID(c context.Context, userID string) (*domain.Profile, error)
}

type profileService struct {
	profileRepository repository.ProfileRepository
	container         container.Container
}

func NewProfileAuth(profileRepository repository.ProfileRepository, container container.Container) ProfileService {
	return &profileService{
		profileRepository: profileRepository,
		container:         container,
	}
}

func (p *profileService) GetProfileByID(ctx context.Context, profileID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, p.container.GetContentTimeout())
	defer cancel()
	return p.profileRepository.GetByID(ctx, profileID)
}
