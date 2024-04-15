package mock

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/stretchr/testify/mock"
)

type ProfileStorageMock struct {
	mock.Mock
}

func (s *ProfileStorageMock) CreateProfile(ctx context.Context, profile *models.Profile) error {
	args := s.Called(profile)
	return args.Error(0)
}

func (s *ProfileStorageMock) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	args := s.Called(id)
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (s *ProfileStorageMock) DeleteProfile(ctx context.Context, id string) error {
	args := s.Called(id)
	return args.Error(0)
}

func (s *ProfileStorageMock) SearchProfile(ctx context.Context, firstName, lastName string) ([]*models.Profile, error) {
	args := s.Called(firstName, lastName)
	return args.Get(0).([]*models.Profile), args.Error(1)
}
