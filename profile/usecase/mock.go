package usecase

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/stretchr/testify/mock"
)

type ProfileUseCaseMock struct {
	mock.Mock
}

func (m *ProfileUseCaseMock) CreateProfile(ctx context.Context, profile *models.Profile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *ProfileUseCaseMock) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileUseCaseMock) DeleteProfile(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProfileUseCaseMock) SearchProfile(ctx context.Context, firstName, lastName string) ([]*models.Profile, error) {
	args := m.Called(firstName, lastName)
	return args.Get(0).([]*models.Profile), args.Error(1)
}
