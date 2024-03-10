package usecase

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth/repository/mock"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	repo := new(mock.LoginStorageMock)
	cfg := &config.Auth{
		HashSalt:   "salt",
		SigningKey: "secret",
		TokenTTL:   86400,
	}
	uc := NewAuthUseCase(repo, cfg)
	var (
		password     = "pass"
		ctx          = context.Background()
		passwordHash = "c8b2505b76926abdc733523caa9f439142f66aa7293a7baaac0aed41a191eef6"
	)

	// Sign Up
	repo.On("CreateLogin", passwordHash).Return("id", nil)
	accountId, err := uc.SignUp(ctx, password)
	assert.NoError(t, err)
	assert.Equal(t, "id", accountId)

	// Sign In (Get Auth Token)
	repo.On("AuthenticateLogin", accountId, passwordHash).Return(nil)
	token, err := uc.SignIn(ctx, accountId, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	//Verify token
	parsedLogin, err := uc.ParseToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, accountId, parsedLogin)
}
