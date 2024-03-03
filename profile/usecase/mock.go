package usecase

import (
	"context"
	"time"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/stretchr/testify/mock"
)

type ProfileUseCaseMock struct {
	mock.Mock
}

func (m *ProfileUseCaseMock) CreateProfile(ctx context.Context, account *models.Login, firstName, lastName string, birthDate time.Time, gender models.Gender, biography, city string) error {
	args := m.Called(account, firstName, lastName, birthDate, gender, biography, city)
	return args.Error(0)
}

func (m *ProfileUseCaseMock) GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error) {
	args := m.Called(account)
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileUseCaseMock) DeleteProfile(ctx context.Context, account *models.Login) error {
	args := m.Called(account)
	return args.Error(0)
}
