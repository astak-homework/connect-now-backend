package usecase

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth/repository/mock"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/astak-homework/connect-now-backend/models"
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
		username = "user"
		password = "pass"
		ctx      = context.Background()
		login    = &models.Login{
			UserName: username,
			Password: "11f5639f22525155cb0b43573ee4212838c78d87", // sha1 of pass+salt
		}
	)

	// Sign Up
	repo.On("CreateLogin", login).Return(nil)
	err := uc.SignUp(ctx, username, password)
	assert.NoError(t, err)

	// Sign In (Get Auth Token)
	repo.On("GetLogin", login.UserName, login.Password).Return(login, nil)
	token, err := uc.SignIn(ctx, username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	//Verify token
	parsedLogin, err := uc.ParseToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, login.UserName, parsedLogin.UserName)
}
