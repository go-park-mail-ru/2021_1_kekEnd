package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/mock"
)

type MoviesUseCaseMock struct {
	mock.Mock
}

func (mockUC *MoviesUseCaseMock) CreateUser(user *models.Movie) error {
	args := mockUC.Called(user)
	return args.Error(0)
}

func (mockUC *MoviesUseCaseMock) GetMovie(movieID string) (*models.Movie, error) {
	args := mockUC.Called(movieID)
	return args.Get(0).(*models.Movie), args.Error(1)
}
