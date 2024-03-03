package profile

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile *models.Profile) error
	GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error)
	DeleteProfile(ctx context.Context, account *models.Login) error
}
