package localcache

import (
	"context"
	"testing"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	id := "id"
	account := &models.Login{ID: id}

	s := NewProfileLocalStorage()

	p := &models.Profile{
		ID:      "id1",
		Account: account,
	}

	err := s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	returnedProfile, err := s.GetProfile(context.Background(), account)
	assert.NoError(t, err)
	assert.Equal(t, p, returnedProfile)
}

func TestDeleteProfile(t *testing.T) {
	id1 := "id1"
	id2 := "id2"

	account1 := &models.Login{ID: id1}
	account2 := &models.Login{ID: id2}

	p := &models.Profile{Account: account1}

	s := NewProfileLocalStorage()

	err := s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	err = s.DeleteProfile(context.Background(), account1)
	assert.NoError(t, err)

	err = s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	err = s.DeleteProfile(context.Background(), account2)
	assert.Error(t, err)
	assert.Equal(t, err, profile.ErrProfileNotFound)
}
