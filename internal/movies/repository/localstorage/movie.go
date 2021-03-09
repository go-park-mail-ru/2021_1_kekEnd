package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"strconv"
	"sync"
)

type MovieLocalStorage struct {
	movies    map[string]*models.Movie
	currentID uint64
	mutex     sync.Mutex
}

func NewMovieLocalStorage() *MovieLocalStorage {
	// dummy data for testing
	movies := map[string]*models.Movie{
		"1": {
			ID:             "1",
			Title:          "Чужой",
			Description:    "Группа космонавтов высаживается на неизвестной планете и знакомится с ксеноморфом. Шедевр Ридли Скотта",
			Voiceover:      []string{"Русский", "Английский"},
			Subtitles:      []string{"Русские"},
			Quality:        "HD",
			ProductionYear: 1979,
			Country:        []string{"Великобритания", "США"},
			Genre:          []string{"Хоррор", "Драма"},
			Slogan:         "slogan",
			Director:       "Ридли Скотт",
			Scriptwriter:   "---",
			Producer:       "---",
			Operator:       "---",
			Composer:       "---",
			Artist:         "---",
			Montage:        "---",
			Budget:         "---",
			Duration:       "---",
			Actors:         []string{"Сигурни Уивер", "Иэн Холм"},
		},
	}

	return &MovieLocalStorage{
		movies:    movies,
		currentID: 2,
	}
}

func (storage *MovieLocalStorage) CreateMovie(movie *models.Movie) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	movie.ID = strconv.FormatUint(storage.currentID, 10)
	storage.movies[movie.ID] = movie
	storage.currentID++

	return nil
}

func (storage *MovieLocalStorage) GetMovieByID(id string) (*models.Movie, error) {
	movie, exists := storage.movies[id]
	if !exists {
		return nil, errors.New("movie not found")
	}
	return movie, nil
}
