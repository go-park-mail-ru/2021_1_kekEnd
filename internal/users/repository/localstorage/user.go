package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

func getHashedPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedPasswordBytes), nil
}

type UserLocalStorage struct {
	users map[string]*models.User
	mutex sync.Mutex
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
		users: users,
	}
}

func (storage *UserLocalStorage) CreateUser(user *models.User) error {
	storage.mutex.Lock()

	hashedPassword, err := getHashedPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

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

func (storage *UserLocalStorage) CheckPassword(password string, user *models.User) (bool, error) {
	hashedPassword, err := getHashedPassword(password)
	if err != nil {
		return false, nil
	}
	return hashedPassword == user.Password, nil
}

func (storage *UserLocalStorage) UpdateUser(username string, newUser *models.User) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	_, exists := storage.users[username]
	if exists {
		storage.users[username] = newUser
		return nil
	}
	return errors.New("user not found")
}
