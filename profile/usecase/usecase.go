package usecase

import (
	"context"
	"time"

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

func (p ProfileUseCase) CreateProfile(ctx context.Context, id, firstName, lastName string, birthDate time.Time, gender models.Gender, biography, city string) error {
	profile := &models.Profile{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Gender:    gender,
		Biography: biography,
		City:      city,
	}
	return p.profileRepo.CreateProfile(ctx, profile)
}

func (p ProfileUseCase) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	return p.profileRepo.GetProfile(ctx, id)
}

func (p ProfileUseCase) DeleteProfile(ctx context.Context, id string) error {
	return p.profileRepo.DeleteProfile(ctx, id)
}
