package localstorage

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/kekEnd_main/models"
	"sync"
)

type UserLocalStorage struct {
	users map[string]*models.User
	mutex *sync.Mutex
}

func NewUserLocalStorage() *UserLocalStorage {
	return &UserLocalStorage{
		users: make(map[string]*models.User),
		mutex: new(sync.Mutex),
	}
}

func (storage *UserLocalStorage) CreateUser(ctx context.Context, user *models.User) error {
	storage.mutex.Lock()
	storage.users[user.ID] = user
	storage.mutex.Unlock()
	return nil
}

func (storage *UserLocalStorage) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for _, user := range storage.users {
		if user.Username == login && user.Password == password {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}
