package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LoginStorageMock struct {
	mock.Mock
}

func (s *LoginStorageMock) CreateLogin(ctx context.Context, password string) (string, error) {
	args := s.Called(password)
	return args.Get(0).(string), args.Error(1)
}

func (s *LoginStorageMock) AuthenticateLogin(ctx context.Context, accountId, password string) error {
	args := s.Called(accountId, password)
	return args.Error(0)
}
