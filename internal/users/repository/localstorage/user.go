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
	users := map[string]*models.User{
		"let_robots_reign": {
			Username:      "let_robots_reign",
			Email:         "sample@ya.ru",
			Password:      "123456789",
			MoviesWatched: 4,
			ReviewsNumber: 2,
		},
	}

	return &UserLocalStorage{
		//users: make(map[string]*models.User),
		users: users,
	}
}

func (storage *UserLocalStorage) CreateUser(user *models.User) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	hashedPassword, err := getHashedPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	storage.users[user.Username] = user
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
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, nil
}

func (storage *UserLocalStorage) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	_, exists := storage.users[user.Username]
	if !exists {
		return nil, errors.New("user not found")
	}

	if user.Username != change.Username {
		return nil, errors.New("username doesn't match")
	}

	if change.Password != "" {
		newPassword, err := getHashedPassword(change.Password)
		if err != nil {
			return nil, err
		}

		user.Password = newPassword
	}

	if change.Email != "" {
		user.Email = change.Email
	}

	if  change.Avatar != "" {
		user.Avatar = change.Avatar
	}

	storage.users[user.Username] = user
	return storage.users[user.Username], nil
}

func (storage *UserLocalStorage) CreateReview(user *models.User, review *models.Review) error {
	// TODO: placeholder until BD appears
	return nil
}
