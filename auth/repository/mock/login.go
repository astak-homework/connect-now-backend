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

func (s *LoginStorageMock) GetPasswordHash(ctx context.Context, accountId string) (string, error) {
	args := s.Called(accountId)
	return args.Get(0).(string), args.Error(1)
}
