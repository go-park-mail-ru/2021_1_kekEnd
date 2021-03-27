package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"strconv"
	"sync"
)

type ReviewLocalStorage struct {
	reviews   map[string]*models.Review
	currentID uint64
	mutex     sync.Mutex
}

func NewReviewLocalStorage() *ReviewLocalStorage {
	return &ReviewLocalStorage{
		reviews: make(map[string]*models.Review),
		currentID: 1,
	}
}

func (storage *ReviewLocalStorage) CreateReview(review *models.Review) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	review.ID = strconv.FormatUint(storage.currentID, 10)
	storage.reviews[review.ID] = review
	storage.currentID++

	return nil
}

func (storage *ReviewLocalStorage) GetUserReviews(username string) []*models.Review {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	var userReviews []*models.Review

	for _, review := range storage.reviews {
		if review.Author == username {
			userReviews = append(userReviews, review)
		}
	}
	return userReviews
}

func (storage *ReviewLocalStorage) GetMovieReviews(movieID string) []*models.Review {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	var movieReviews []*models.Review

	for _, review := range storage.reviews {
		if review.MovieID == movieID {
			movieReviews = append(movieReviews, review)
		}
	}
	return movieReviews
}

func (storage *ReviewLocalStorage) GetUserReviewForMovie(username string, movieID string) (*models.Review, error)  {
	userReviews := storage.GetUserReviews(username)
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for _, review := range userReviews {
		if review.MovieID == movieID {
			return review, nil
		}
	}
	return nil, errors.New("review doesn't exist")
}

func (storage *ReviewLocalStorage) EditUserReviewForMovie(review *models.Review) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	storage.reviews[review.ID] = review
	return nil
}рецензии

func (storage *ReviewLocalStorage) DeleteUserReviewForMovie(username string, movieID string) error {
	userReview, err := storage.GetUserReviewForMovie(username, movieID)
	if err != nil {
		return err
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	delete(storage.reviews, userReview.ID)
	return nil
}
