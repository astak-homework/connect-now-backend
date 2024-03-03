package localstorage

import (
	"context"
	"sync"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/google/uuid"
)

type LoginLocalStorage struct {
	logins map[string]*models.Login
	mutex  *sync.Mutex
}

func NewLoginLocalStorage() *LoginLocalStorage {
	return &LoginLocalStorage{
		logins: make(map[string]*models.Login),
		mutex:  new(sync.Mutex),
	}
}

func (s *LoginLocalStorage) CreateLogin(ctx context.Context, login *models.Login) error {
	s.mutex.Lock()
	login.ID = uuid.NewString()
	s.logins[login.ID] = login
	s.mutex.Unlock()
	return nil
}

func (s *LoginLocalStorage) GetLogin(ctx context.Context, username, password string) (*models.Login, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, login := range s.logins {
		if login.UserName == username && login.Password == password {
			return login, nil
		}
	}

	return nil, auth.ErrUserNotFound
}
