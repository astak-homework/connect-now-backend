package localstorage

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	s := NewLoginLocalStorage()

	id1 := "id"

	login := &models.Login{
		ID:       id1,
		UserName: "user",
		Password: "password",
	}

	err := s.CreateLogin(context.Background(), login)
	assert.NoError(t, err)

	returnedLogin, err := s.GetLogin(context.Background(), "user", "password")
	assert.NoError(t, err)
	assert.Equal(t, login, returnedLogin)

	_, err = s.GetLogin(context.Background(), "user", "")
	assert.Error(t, err)
	assert.Equal(t, err, auth.ErrUserNotFound)
}
