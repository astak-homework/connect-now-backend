package profile

import (
	"context"
	"time"

	"github.com/astak-homework/connect-now-backend/models"
)

type UseCase interface {
	CreateProfile(ctx context.Context, account *models.Login, firstName, lastName string, birthDate time.Time, gender models.Gender, biography, city string) error
	GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error)
	DeleteProfile(ctx context.Context, account *models.Login) error
}
