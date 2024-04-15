package profile

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
)

type UseCase interface {
	CreateProfile(ctx context.Context, profile *models.Profile) error
	GetProfile(ctx context.Context, id string) (*models.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
	SearchProfile(ctx context.Context, firstName, lastName string) ([]*models.Profile, error)
}
