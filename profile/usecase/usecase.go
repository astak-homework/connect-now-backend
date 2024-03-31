package usecase

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
)

type ProfileUseCase struct {
	profileRepo profile.Repository
}

func NewProfileUseCase(profileRepo profile.Repository) *ProfileUseCase {
	return &ProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (p ProfileUseCase) CreateProfile(ctx context.Context, profile *models.Profile) error {
	return p.profileRepo.CreateProfile(ctx, profile)
}

func (p ProfileUseCase) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	return p.profileRepo.GetProfile(ctx, id)
}

func (p ProfileUseCase) DeleteProfile(ctx context.Context, id string) error {
	return p.profileRepo.DeleteProfile(ctx, id)
}

func (p ProfileUseCase) SearchProfile(ctx context.Context, firstName, lastName string) ([]*models.Profile, error) {
	return p.profileRepo.SearchProfile(ctx, firstName, lastName)
}
