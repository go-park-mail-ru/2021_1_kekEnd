package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"sync"
)

type UserLocalStorage struct {
	users   map[string]*models.User
	mutex   *sync.Mutex
}

func NewUserLocalStorage() *UserLocalStorage {
	// dummy data for testing
	users := map[string]*models.User{
		"1": {
			Username:      "let-robots-reign",
			Email:         "sample@ya.ru",
			Password:      "1234",
			MoviesWatched: 4,
			ReviewsNumber: 2,
		},
	}

	return &UserLocalStorage{
		users:   users,
		mutex:   new(sync.Mutex),
	}
}

func (storage *UserLocalStorage) CreateUser(user *models.User) error {
	storage.mutex.Lock()

	storage.users[user.Username] = user

	storage.mutex.Unlock()
	return nil
}

func (storage *UserLocalStorage) GetUserByUsername(username string) (*models.User, error) {
	user, exists := storage.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (storage *UserLocalStorage) UpdateUser(username string, newUser *models.User) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if _, exists := storage.users[username]; exists {
		storage.users[username] = newUser
		return nil
	}
	return storage.CreateUser(newUser)
}
