package usecase

import (
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

// UsersUseCase структура usecase юзера
type UsersUseCase struct {
	userRepository    users.UserRepository
	reviewsRepository reviews.ReviewRepository
	ratingsRepository ratings.Repository
	actorsRepository  actors.Repository
}

// NewUsersUseCase инициализация usecase юзера
func NewUsersUseCase(repo users.UserRepository, reviews reviews.ReviewRepository, ratings ratings.Repository,
	actors actors.Repository) *UsersUseCase {
	return &UsersUseCase{
		userRepository:    repo,
		reviewsRepository: reviews,
		ratingsRepository: ratings,
		actorsRepository:  actors,
	}
}

// CreateUser создание юзера
func (usersUC *UsersUseCase) CreateUser(user *models.User) error {
	_, err := usersUC.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	return usersUC.userRepository.CreateUser(user)
}

// Login логин юзера
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

// GetUser получить юзера
func (usersUC *UsersUseCase) GetUser(username string) (*models.User, error) {
	user, err := usersUC.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	favActors, err := usersUC.actorsRepository.GetFavoriteActors(user.Username)
	if err != nil {
		return nil, err
	}
	user.FavoriteActors = favActors
	return user, nil
}

// UpdateUser обновить юзера
func (usersUC *UsersUseCase) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	err := usersUC.userRepository.CheckEmailUnique(change.Email)
	if err != nil {
		return nil, err
	}

	return usersUC.userRepository.UpdateUser(user, change)
}

func (usersUC *UsersUseCase) toggleSubscribe(subscriber string, user string, isSubscribing bool) error {
	var err error
	if isSubscribing {
		err = usersUC.userRepository.Subscribe(subscriber, user)
	} else {
		err = usersUC.userRepository.Unsubscribe(subscriber, user)
	}
	if err != nil {
		return err
	}

	subscriberModel, err := usersUC.GetUser(subscriber)
	if err != nil {
		return err
	}
	var newSubscriptionsNum uint
	if isSubscribing {
		newSubscriptionsNum = *subscriberModel.Subscriptions + 1
	} else {
		newSubscriptionsNum = *subscriberModel.Subscriptions - 1
	}
	_, err = usersUC.UpdateUser(subscriberModel, models.User{
		Username:      subscriber,
		Subscriptions: &newSubscriptionsNum,
	})
	if err != nil {
		return err
	}

	userModel, err := usersUC.GetUser(user)
	if err != nil {
		return err
	}
	var newSubscribersNum uint
	if isSubscribing {
		newSubscribersNum = *userModel.Subscribers + 1
	} else {
		newSubscribersNum = *userModel.Subscribers - 1
	}

	_, err = usersUC.UpdateUser(userModel, models.User{
		Username:    user,
		Subscribers: &newSubscribersNum,
	})

	if err != nil {
		return err
	}
	return nil
}

// Subscribe подписаться на юзера
func (usersUC *UsersUseCase) Subscribe(subscriber string, user string) error {
	unsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)
	if err != nil {
		return err
	}
	if unsubscribed {
		return usersUC.toggleSubscribe(subscriber, user, true)
	}
	return fmt.Errorf("%s is already subscribed to %s", subscriber, user)
}

// Unsubscribe отписаться от юзера
func (usersUC *UsersUseCase) Unsubscribe(subscriber string, user string) error {
	unsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)
	if err != nil {
		return err
	}
	if !unsubscribed {
		return usersUC.toggleSubscribe(subscriber, user, false)
	}
	return fmt.Errorf("%s is not subscribed to %s", subscriber, user)
}

// GetSubscribers получить подписчиков
func (usersUC *UsersUseCase) GetSubscribers(page int, user string) (int, []models.UserNoPassword, error) {
	startIndex := (page - 1) * constants.SubsPageSize
	return usersUC.userRepository.GetSubscribers(startIndex, user)
}

// IsSubscribed проверить подписан ли
func (usersUC *UsersUseCase) IsSubscribed(subscriber string, user string) (bool, error) {
	isUnsubscribed, err := usersUC.userRepository.CheckUnsubscribed(subscriber, user)
	return !isUnsubscribed, err
}

// GetSubscriptions получить подписки
func (usersUC *UsersUseCase) GetSubscriptions(page int, user string) (int, []models.UserNoPassword, error) {
	startIndex := (page - 1) * constants.SubsPageSize
	return usersUC.userRepository.GetSubscriptions(startIndex, user)
}

// GetFeed получить новости
func (usersUC *UsersUseCase) GetFeed(username string) (models.Feed, error) {
	_, subs, err := usersUC.userRepository.GetSubscriptions(0, username)
	if err != nil {
		return models.Feed{}, err
	}

	reviewsFeed, err := usersUC.reviewsRepository.GetFeed(subs)
	if err != nil {
		return models.Feed{}, err
	}

	ratingsFeed, err := usersUC.ratingsRepository.GetFeed(subs)
	if err != nil {
		return models.Feed{}, err
	}

	feed := models.Feed{
		Ratings: ratingsFeed,
		Reviews: reviewsFeed,
	}

	return feed, nil
}
