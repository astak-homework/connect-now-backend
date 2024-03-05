package profile

import (
	"context"
	"time"

	"github.com/astak-homework/connect-now-backend/models"
)

type UseCase interface {
	CreateProfile(ctx context.Context, id, firstName, lastName string, birthDate time.Time, gender models.Gender, biography, city string) error
	GetProfile(ctx context.Context, id string) (*models.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
}
