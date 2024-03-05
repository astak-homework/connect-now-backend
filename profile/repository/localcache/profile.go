package localcache

import (
	"context"
	"sync"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
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
	s.profiles[profile.ID] = profile
	s.mutex.Unlock()
	return nil
}

func (s *ProfileLocalStorage) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	p, ok := s.profiles[id]
	if !ok {
		return nil, profile.ErrProfileNotFound
	}

	return p, nil
}

func (s *ProfileLocalStorage) DeleteProfile(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	p, ok := s.profiles[id]
	if !ok {
		return profile.ErrProfileNotFound
	}

	delete(s.profiles, p.ID)
	return nil
}
