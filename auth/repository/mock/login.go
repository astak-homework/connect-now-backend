package mock

import (
	"context"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/stretchr/testify/mock"
)

type LoginStorageMock struct {
	mock.Mock
}

func (s *LoginStorageMock) CreateLogin(ctx context.Context, login *models.Login) error {
	args := s.Called(login)
	return args.Error(0)
}

func (s *LoginStorageMock) GetLogin(ctx context.Context, username, password string) (*models.Login, error) {
	args := s.Called(username, password)
	return args.Get(0).(*models.Login), args.Error(1)
}
