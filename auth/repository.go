package auth

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
)

type LoginRepository interface {
	CreateLogin(ctx context.Context, login *models.Login) error
	GetLogin(ctx context.Context, username, password string) (*models.Login, error)
}
