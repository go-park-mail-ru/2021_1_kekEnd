package repository

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/mock"
)

type MoviesStorageMock struct {
	mock.Mock
}

func (storageMock *MoviesStorageMock) CreateMovie(movie *models.Movie) error {
	args := storageMock.Called(movie)
	return args.Error(0)
}

func (storageMock *MoviesStorageMock) GetMovieByID(movieID string) (*models.Movie, error) {
	args := storageMock.Called(movieID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), args.Error(1)
}
