package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"sync"
)

type UserLocalStorage struct {
	users   map[int]*models.User
	counter int
	mutex   *sync.Mutex
}

func NewUserLocalStorage() *UserLocalStorage {
	return &UserLocalStorage{
		users:   make(map[int]*models.User),
		counter: 1,
		mutex:   new(sync.Mutex),
	}
}

func (storage *UserLocalStorage) CreateUser(user *models.User) error {
	storage.mutex.Lock()

	user.ID = storage.counter
	storage.users[user.ID] = user
	storage.counter++

	storage.mutex.Unlock()
	return nil
}

func (storage *UserLocalStorage) GetUserByLoginPassword(login, password string) (*models.User, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for _, user := range storage.users {
		if user.Username == login && user.Password == password {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (storage *UserLocalStorage) GetUserByID(id int) (*models.User, error) {
	user, exists := storage.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (storage *UserLocalStorage) UpdateUser(id int, newUser *models.User) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if _, exists := storage.users[id]; exists {
		storage.users[id] = newUser
		return nil
	}
	return storage.CreateUser(newUser)
}
