package usecase

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth/repository/mock"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthFlow(t *testing.T) {
	repo := new(mock.LoginStorageMock)
	cfg := &config.Auth{
		SigningKey: "secret",
		TokenTTL:   86400,
	}
	uc := NewAuthUseCase(repo, cfg)
	var (
		password = "pass"
		ctx      = context.Background()
	)

	// Sign Up
	var passwordHash string
	repo.On("CreateLogin", testifyMock.AnythingOfType("string")).Return("id", nil).Run(func(args testifyMock.Arguments) {
		passwordHash = args.Get(0).(string)
	})
	accountId, err := uc.SignUp(ctx, password)
	assert.NoError(t, err)
	assert.Equal(t, "id", accountId)

	// Sign In (Get Auth Token)
	repo.On("GetPasswordHash", accountId).Return(passwordHash, nil)
	token, err := uc.SignIn(ctx, accountId, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	//Verify token
	parsedLogin, err := uc.ParseToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, accountId, parsedLogin)
}
