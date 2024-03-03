package localcache

import (
	"context"
	"sync"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/google/uuid"
)

type ProfileLocalStorage struct {
	profiles map[string]*models.Profile
	mutex    *sync.Mutex
}

func NewProfileLocalStorage() *ProfileLocalStorage {
	return &ProfileLocalStorage{
		profiles: make(map[string]*models.Profile),
		mutex:    new(sync.Mutex),
	}
}

func (s *ProfileLocalStorage) CreateProfile(ctx context.Context, profile *models.Profile) error {
	s.mutex.Lock()
	profile.ID = uuid.NewString()
	s.profiles[profile.ID] = profile
	s.mutex.Unlock()
	return nil
}

func (s *ProfileLocalStorage) GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, p := range s.profiles {
		if p.Account.ID == account.ID {
			return p, nil
		}
	}

	return nil, profile.ErrProfileNotFound
}

func (s *ProfileLocalStorage) DeleteProfile(ctx context.Context, account *models.Login) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ok := false
	for key, p := range s.profiles {
		if p.Account.ID == account.ID {
			delete(s.profiles, key)
			ok = true
		}
	}

	if !ok {
		return profile.ErrProfileNotFound
	}

	return nil
}
