package auth

import (
	"context"
)

const CtxLoginKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, password string) (string, error)
	SignIn(ctx context.Context, accountId, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (string, error)
}
