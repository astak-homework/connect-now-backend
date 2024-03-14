package localstorage

import (
	"context"
	"sync"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/google/uuid"
)

type LoginLocalStorage struct {
	logins map[string]string
	mutex  *sync.Mutex
}

func NewLoginLocalStorage() *LoginLocalStorage {
	return &LoginLocalStorage{
		logins: make(map[string]string),
		mutex:  new(sync.Mutex),
	}
}

func (s *LoginLocalStorage) CreateLogin(ctx context.Context, password string) (string, error) {
	s.mutex.Lock()
	accountId := uuid.NewString()
	s.logins[accountId] = password
	s.mutex.Unlock()
	return accountId, nil
}

func (s *LoginLocalStorage) GetPasswordHash(ctx context.Context, accountId string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	h, ok := s.logins[accountId]
	if !ok {
		return "", auth.ErrUserNotFound
	}

	return h, nil
}
