package localstorage

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	s := NewLoginLocalStorage()

	accountId, err := s.CreateLogin(context.Background(), "password")
	assert.NoError(t, err)

	password, err := s.GetPasswordHash(context.Background(), accountId)
	assert.NoError(t, err)
	assert.Equal(t, "password", password)

	_, err = s.GetPasswordHash(context.Background(), "12345")
	assert.ErrorIs(t, auth.ErrUserNotFound, err)
}
