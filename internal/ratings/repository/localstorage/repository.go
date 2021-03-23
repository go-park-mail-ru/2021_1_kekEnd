package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"sync"
)

type RatingsLocalStorage struct {
	ratings   []models.Rating
	currentID int
	mutex     sync.Mutex
}

func NewRatingsLocalStorage() *RatingsLocalStorage {
	ratings := []models.Rating{
		{
			UserID:  "let_robots_reign",
			MovieID: "1",
			Score:   8,
		},
		{
			UserID:  "let_robots_reign",
			MovieID: "4",
			Score:   9,
		},
		{
			UserID:  "let_robots_reign",
			MovieID: "5",
			Score:   10,
		},
	}

	return &RatingsLocalStorage{
		ratings:   ratings,
		currentID: 3,
	}
}

func (r *RatingsLocalStorage) CreateRating(userID string, movieID string, score int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	Rating := models.Rating{
		UserID:  userID,
		MovieID: movieID,
		Score:   score,
	}

	r.ratings = append(r.ratings, Rating)
	r.currentID++

	return nil
}

func (r *RatingsLocalStorage) GetRating(userID string, movieID string) (models.Rating, error) {
	index := -1
	for i, s := range r.ratings {
		if s.UserID == userID && s.MovieID == movieID {
			index = i
			break
		}
	}

	if index != -1 {
		return r.ratings[index], nil
	}

	return models.Rating{}, errors.New("rating doesn't exist")
}

func (r *RatingsLocalStorage) DeleteRating(userID string, movieID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.currentID == 0 {
		return errors.New("no ratings in database")
	}

	deleted := false
	for i, s := range r.ratings {
		if s.UserID == userID && s.MovieID == movieID {
			r.ratings[i] = r.ratings[r.currentID-1]
			r.ratings = r.ratings[:r.currentID-1]
			r.currentID--
			deleted = true
			break
		}
	}
	if !deleted {
		return errors.New("rating doesn't exist")
	}
	return nil
}

func (r *RatingsLocalStorage) UpdateRating(userID string, movieID string, score int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, s := range r.ratings {
		if s.UserID == userID && s.MovieID == movieID {
			s.Score = score
			return nil
		}
	}

	return errors.New("rating doesn't exist")
}
