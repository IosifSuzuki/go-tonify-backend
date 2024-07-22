package usecase

import (
	"context"
	"go-tonify-backend/internal/domain"
	"go-tonify-backend/internal/repository"
	"time"
)

type ProfileUseCase interface {
	GetProfileByID(c context.Context, userID string) (*domain.Profile, error)
}

type profileUseCase struct {
	profileRepository repository.ProfileRepository
	contextTimeout    time.Duration
}

func NewProfileUseCase(profileRepository repository.ProfileRepository, timeout time.Duration) ProfileUseCase {
	return &profileUseCase{
		profileRepository: profileRepository,
		contextTimeout:    timeout,
	}
}

func (p *profileUseCase) GetProfileByID(ctx context.Context, profileID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()
	return p.profileRepository.GetByID(ctx, profileID)
}
