package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) SignUp(ctx context.Context, password string) (string, error) {
	args := m.Called(password)
	return args.Get(0).(string), args.Error(1)
}

func (m *AuthUseCaseMock) SignIn(ctx context.Context, accountId, password string) (string, error) {
	args := m.Called(accountId, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *AuthUseCaseMock) ParseToken(ctx context.Context, accessToken string) (string, error) {
	args := m.Called(accessToken)
	return args.Get(0).(string), args.Error(1)
}
