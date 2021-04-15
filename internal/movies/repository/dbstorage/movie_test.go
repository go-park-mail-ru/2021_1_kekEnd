package localstorage

import (
	"context"

	// "fmt"
	"testing"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	// "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/pashagolub/pgxmock"
)


func TestGetMovieByID(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewMovieRepository(mock)
	movie := &models.Movie{
		ID: "1",
		Title: "Good film",
		Description: "Veryvery well",
		ProductionYear: 1999,
		Country: []string{"Russia"},
		Genre: []string{"comedy"},
		Slogan: "BEST",
		Director: "Luk Besson",
		Scriptwriter: "ilya228",
		Producer: "Luk Besson",
		Operator: "Luk Besson",
		Composer: "ilya228",
		Artist: "ilya228",
		Montage: "ilya228",
		Budget: "ilya228",
		Duration: "ilya228",
		Actors: []models.ActorData{{1, "ilya"}},
		Poster: "qwe",
		Banner: "qwe",
		TrailerPreview: "qwe",
		Rating: 1,
		RatingCount: 1,
	}

	rows := pgxmock.NewRows([]string{"id", "title", "description", "productionYear", "country",
               "genre", "slogan", "director", "scriptwriter", "producer", "operator", "composer",
               "artist", "montage", "budget", "duration", "poster", "banner", "trailerPreview",
               "ROUND(CAST(rating AS numeric), 1) AS rating", "rating_count"}).
	AddRow(1,
		movie.Title, movie.Description,
		movie.ProductionYear, movie.Country, movie.Genre, movie.Slogan, movie.Director,
		movie.Scriptwriter, movie.Producer, movie.Operator, movie.Composer, movie.Artist,
		movie.Montage, movie.Budget, movie.Duration,  movie.Poster,
		movie.Banner, movie.TrailerPreview, movie.Rating, movie.RatingCount)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(movie.ID).WillReturnRows(rows)
	mock.ExpectCommit()

	if _, err = movieRepo.GetMovieByID(movie.ID); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBestMovies(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewMovieRepository(mock)
	movie := &models.Movie{
		ID: "1",
		Title: "Good film",
		Description: "Veryvery well",
		ProductionYear: 1999,
		Country: []string{"Russia"},
		Genre: []string{"comedy"},
		Slogan: "BEST",
		Director: "Luk Besson",
		Scriptwriter: "ilya228",
		Producer: "Luk Besson",
		Operator: "Luk Besson",
		Composer: "ilya228",
		Artist: "ilya228",
		Montage: "ilya228",
		Budget: "ilya228",
		Duration: "ilya228",
		Actors: []models.ActorData{{1, "ilya"}},
		Poster: "qwe",
		Banner: "qwe",
		TrailerPreview: "qwe",
		Rating: 1,
		RatingCount: 1,
	}

	rows := pgxmock.NewRows([]string{"id", "title", "description", "productionYear", "country",
               "genre", "slogan", "director", "scriptwriter", "producer", "operator", "composer",
               "artist", "montage", "budget", "duration", "poster", "banner", "trailerPreview",
               "ROUND(CAST(rating AS numeric), 1) AS rating", "rating_count"}).
	AddRow(1,
		movie.Title, movie.Description,
		movie.ProductionYear, movie.Country, movie.Genre, movie.Slogan, movie.Director,
		movie.Scriptwriter, movie.Producer, movie.Operator, movie.Composer, movie.Artist,
		movie.Montage, movie.Budget, movie.Duration,  movie.Poster,
		movie.Banner, movie.TrailerPreview, movie.Rating, movie.RatingCount)

	rows2 := pgxmock.NewRows([]string{"movies_count"}).AddRow(0)

	mock.ExpectQuery("SELECT").WillReturnRows(rows2)

	mock.ExpectQuery("SELECT").WithArgs(movie.ID).WillReturnRows(rows)

	if _, _, err = movieRepo.GetBestMovies(1); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllGenres(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewMovieRepository(mock)

	rows := pgxmock.NewRows([]string{"available_genres"}).AddRow("comedy")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	if _, err = movieRepo.GetAllGenres(); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMoviesByGenres(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewMovieRepository(mock)
	movie := &models.Movie{
		ID: "1",
		Title: "Good film",
		Description: "Veryvery well",
		ProductionYear: 1999,
		Country: []string{"Russia"},
		Genre: []string{"comedy"},
		Slogan: "BEST",
		Director: "Luk Besson",
		Scriptwriter: "ilya228",
		Producer: "Luk Besson",
		Operator: "Luk Besson",
		Composer: "ilya228",
		Artist: "ilya228",
		Montage: "ilya228",
		Budget: "ilya228",
		Duration: "ilya228",
		Actors: []models.ActorData{{1, "ilya"}},
		Poster: "qwe",
		Banner: "qwe",
		TrailerPreview: "qwe",
		Rating: 1,
		RatingCount: 1,
	}

	rows := pgxmock.NewRows([]string{"id", "title", "description", "productionYear", "country",
               "genre", "slogan", "director", "scriptwriter", "producer", "operator", "composer",
               "artist", "montage", "budget", "duration", "poster", "banner", "trailerPreview",
               "ROUND(CAST(rating AS numeric), 1) AS rating", "rating_count"}).
	AddRow(1,
		movie.Title, movie.Description,
		movie.ProductionYear, movie.Country, movie.Genre, movie.Slogan, movie.Director,
		movie.Scriptwriter, movie.Producer, movie.Operator, movie.Composer, movie.Artist,
		movie.Montage, movie.Budget, movie.Duration,  movie.Poster,
		movie.Banner, movie.TrailerPreview, movie.Rating, movie.RatingCount)

	rows2 := pgxmock.NewRows([]string{"movies_count"}).AddRow(0)

	mock.ExpectQuery("SELECT").WillReturnRows(rows2)

	mock.ExpectQuery("SELECT").WithArgs(movie.ID).WillReturnRows(rows)

	// now we execute our method
	if _, _, err = movieRepo.GetMoviesByGenres([]string{"comedy"}, 1); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetActorsData(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewMovieRepository(mock)

	rows := pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "ilya").AddRow(2, "qwe")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	if _, err = movieRepo.getActorsData([]string{"ilya", "qwe"}); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}