package usecase

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

type UsersUseCase struct {
	userRepository    users.UserRepository
	reviewsRepository reviews.ReviewRepository
	actorsRepository  actors.Repository
}

func NewUsersUseCase(repo users.UserRepository, reviews reviews.ReviewRepository, actors actors.Repository) *UsersUseCase {
	return &UsersUseCase{
		userRepository:    repo,
		reviewsRepository: reviews,
		actorsRepository:  actors,
	}
}

func (usersUC *UsersUseCase) CreateUser(user *models.User) error {
	_, err := usersUC.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	return usersUC.userRepository.CreateUser(user)
}

func (usersUC *UsersUseCase) Login(login, password string) bool {
	user, err := usersUC.userRepository.GetUserByUsername(login)
	if err != nil {
		return false
	}
	correct, err := usersUC.userRepository.CheckPassword(password, user)
	if err != nil {
		return false
	}
	return correct
}

func (usersUC *UsersUseCase) GetUser(username string) (*models.User, error) {
	user, err := usersUC.userRepository.GetUserByUsername(username)
	if err != nil {
		return &models.User{}, err
	}
	favActors, err := usersUC.actorsRepository.GetFavoriteActors(user.Username)
	if err != nil {
		return &models.User{}, err
	}
	user.FavoriteActors = favActors
	return user, nil
}

func (usersUC *UsersUseCase) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	err := usersUC.userRepository.CheckEmailUnique(change.Email)
	if err != nil {
		return nil, err
	}

	return usersUC.userRepository.UpdateUser(user, change)
}

func (usersUC *UsersUseCase) Subscribe(subscriber string, user string) error {
	unsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)

	if err != nil {
		return err
	}

	if unsubscribed {
		return usersUC.userRepository.Subscribe(subscriber, user)
	}

	return fmt.Errorf("%s is already subscribed to %s", subscriber, user)
}

func (usersUC *UsersUseCase) Unsubscribe(subscriber string, user string) error {
	unsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)

	if err != nil {
		return err
	}

	if !unsubscribed {
		return usersUC.userRepository.Unsubscribe(subscriber, user)
	}

	return fmt.Errorf("%s is not subscribed to %s", subscriber, user)
}

func (usersUC *UsersUseCase) GetSubscribers(page int, user string) (int, []*models.UserNoPassword, error) {
	startIndex := (page - 1) * _const.SubsPageSize
	return usersUC.userRepository.GetSubscribers(startIndex, user)
}

func (usersUC *UsersUseCase) IsSubscribed(subscriber string, user string) (bool, error) {
	isUnsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)
	return !isUnsubscribed, err
}

func (usersUC *UsersUseCase) GetSubscriptions(page int, user string) (int, []*models.UserNoPassword, error) {
	startIndex := (page - 1) * _const.SubsPageSize
	return usersUC.userRepository.GetSubscriptions(startIndex, user)
}

func (usersUC *UsersUseCase) GetFeed(username string) ([]*models.Notification, error) {
	_, subs, err := usersUC.userRepository.GetSubscriptions(0, username)
	if err != nil {
		return nil, err
	}

	return usersUC.reviewsRepository.GetFeed(subs)

}
