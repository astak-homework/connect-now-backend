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

func (p ProfileUseCase) CreateProfile(ctx context.Context, account *models.Login, firstName, lastName string, birthDate time.Time, gender models.Gender, biography, city string) error {
	profile := &models.Profile{
		Account:   account,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Gender:    gender,
		Biography: biography,
		City:      city,
	}
	return p.profileRepo.CreateProfile(ctx, profile)
}

func (p ProfileUseCase) GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error) {
	return p.profileRepo.GetProfile(ctx, account)
}

func (p ProfileUseCase) DeleteProfile(ctx context.Context, account *models.Login) error {
	return p.profileRepo.DeleteProfile(ctx, account)
}
