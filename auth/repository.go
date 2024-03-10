package auth

import (
	"context"
)

type LoginRepository interface {
	CreateLogin(ctx context.Context, passwordHash string) (string, error)
	AuthenticateLogin(ctx context.Context, accountId, password_hash string) error
}
