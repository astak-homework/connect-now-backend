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
	s := NewProfileLocalStorage()

	p := &models.Profile{
		ID: id,
	}

	err := s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	returnedProfile, err := s.GetProfile(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, p, returnedProfile)
}

func TestDeleteProfile(t *testing.T) {
	id1 := "id1"
	id2 := "id2"

	p := &models.Profile{ID: id1}

	s := NewProfileLocalStorage()

	err := s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	err = s.DeleteProfile(context.Background(), id1)
	assert.NoError(t, err)

	err = s.CreateProfile(context.Background(), p)
	assert.NoError(t, err)

	err = s.DeleteProfile(context.Background(), id2)
	assert.Error(t, err)
	assert.Equal(t, err, profile.ErrProfileNotFound)
}

func TestSearchProfile(t *testing.T) {
	s := NewProfileLocalStorage()

	p1 := &models.Profile{
		ID:        "p1",
		FirstName: "Мариям",
		LastName:  "Петрушевская",
	}
	p2 := &models.Profile{
		ID:        "p2",
		FirstName: "Альвиан",
		LastName:  "Прокопов",
	}

	err := s.CreateProfile(context.Background(), p1)
	assert.NoError(t, err)

	err = s.CreateProfile(context.Background(), p2)
	assert.NoError(t, err)

	profiles, err := s.SearchProfile(context.Background(), "М", "П")
	assert.NoError(t, err)

	assert.Len(t, profiles, 1)
	assert.Equal(t, p1.ID, profiles[0].ID)
}
