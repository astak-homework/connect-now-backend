package auth

import (
	"context"
)

type LoginRepository interface {
	CreateLogin(ctx context.Context, passwordHash string) (string, error)
	GetPasswordHash(ctx context.Context, accountId string) (string, error)
}
