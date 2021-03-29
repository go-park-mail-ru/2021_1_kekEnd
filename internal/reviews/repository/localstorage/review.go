package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"math"
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
		// for pagination testing
		reviews: map[string]*models.Review{
			"1": {
				ID:         "1",
				Title:      "1",
				ReviewType: "positive",
				Content:    "1",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"2": {
				ID:         "2",
				Title:      "2",
				ReviewType: "positive",
				Content:    "2",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"3": {
				ID:         "3",
				Title:      "3",
				ReviewType: "positive",
				Content:    "3",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"4": {
				ID:         "4",
				Title:      "4",
				ReviewType: "positive",
				Content:    "4",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"5": {
				ID:         "5",
				Title:      "5",
				ReviewType: "positive",
				Content:    "5",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"6": {
				ID:         "6",
				Title:      "6",
				ReviewType: "positive",
				Content:    "6",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"7": {
				ID:         "7",
				Title:      "7",
				ReviewType: "positive",
				Content:    "7",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"8": {
				ID:         "8",
				Title:      "8",
				ReviewType: "positive",
				Content:    "8",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"9": {
				ID:         "9",
				Title:      "9",
				ReviewType: "positive",
				Content:    "9",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"10": {
				ID:         "10",
				Title:      "10",
				ReviewType: "positive",
				Content:    "10",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"11": {
				ID:         "11",
				Title:      "11",
				ReviewType: "positive",
				Content:    "11",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"12": {
				ID:         "12",
				Title:      "12",
				ReviewType: "positive",
				Content:    "12",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"13": {
				ID:         "13",
				Title:      "13",
				ReviewType: "positive",
				Content:    "13",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"14": {
				ID:         "14",
				Title:      "14",
				ReviewType: "positive",
				Content:    "14",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
			"15": {
				ID:         "15",
				Title:      "15",
				ReviewType: "positive",
				Content:    "15",
				Author:     "let_robots_reign",
				MovieID:    "1",
			},
		},
		currentID: 16,
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

func (storage *ReviewLocalStorage) GetMovieReviews(movieID string, page int) (int, []*models.Review) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	var movieReviews []*models.Review

	for _, review := range storage.reviews {
		if review.MovieID == movieID {
			movieReviews = append(movieReviews, review)
		}
	}

	startIndex := (page - 1) * _const.ReviewsPageSize
	endIndex := startIndex + _const.ReviewsPageSize
	pagesNumber := int(math.Ceil(float64(len(storage.reviews)) / _const.ReviewsPageSize))

	return pagesNumber, movieReviews[startIndex:endIndex]
}

func (storage *ReviewLocalStorage) GetUserReviewForMovie(username string, movieID string) (*models.Review, error) {
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
}

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
