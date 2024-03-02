package auth

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
)

const CtxLoginKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.Login, error)
}
