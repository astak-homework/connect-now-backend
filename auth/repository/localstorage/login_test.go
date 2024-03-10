package localstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	s := NewLoginLocalStorage()

	accountId, err := s.CreateLogin(context.Background(), "password")
	assert.NoError(t, err)

	err = s.AuthenticateLogin(context.Background(), accountId, "password")
	assert.NoError(t, err)

	err = s.AuthenticateLogin(context.Background(), accountId, "")
	assert.Error(t, err)
}
