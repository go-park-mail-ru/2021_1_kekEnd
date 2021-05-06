package dbstorage

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/pashagolub/pgxmock"
)

func TestGetActorByID(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewActorRepository(mock)
	actor := &models.Actor{
		ID:           "1",
		Name:         "Ivan",
		Biography:    "from russia",
		BirthDate:    "1.1.1990",
		Origin:       "russia",
		Profession:   "actor",
		MoviesCount:  5,
		MoviesRating: 5,
		Movies:       []models.MovieReference{{"1", "QWE", 7.5}},
		Avatar:       "qwe",
	}

	rows := pgxmock.NewRows([]string{"id", "name", "biography", "birthdate", "origin", "profession", "avatar"}).
		AddRow(1, actor.Name, actor.Biography, actor.BirthDate, actor.Origin, actor.Profession, actor.Avatar)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	rows2 := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT").WithArgs(actor.Name, 1).WillReturnRows(rows2)

	rows3 := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT").WithArgs(actor.ID).WillReturnRows(rows3)

	rows4 := pgxmock.NewRows([]string{"id", "title", "rnd"}).
		AddRow(1, actor.Movies[0].Title, actor.Movies[0].Rating)
	mock.ExpectQuery("SELECT").WithArgs(actor.ID, 10).WillReturnRows(rows4)

	if _, err = movieRepo.GetActorByID(actor.ID, actor.Name); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMoviesForActor(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewActorRepository(mock)
	actor := &models.Actor{
		ID:           "1",
		Name:         "Ivan",
		Biography:    "from russia",
		BirthDate:    "1.1.1990",
		Origin:       "russia",
		Profession:   "actor",
		MoviesCount:  5,
		MoviesRating: 5,
		Movies:       []models.MovieReference{{"1", "QWE", 7.5}},
		Avatar:       "qwe",
	}

	rows := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT").WithArgs(actor.ID).WillReturnRows(rows)

	rows2 := pgxmock.NewRows([]string{"id", "title", "rnd"}).
		AddRow(1, actor.Movies[0].Title, actor.Movies[0].Rating)
	mock.ExpectQuery("SELECT").WithArgs(actor.ID, 10).WillReturnRows(rows2)

	if _, _, err = movieRepo.getMoviesForActor(actor.ID); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetFavoriteActors(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewActorRepository(mock)
	actor := &models.Actor{
		ID:           "1",
		Name:         "Ivan",
		Biography:    "from russia",
		BirthDate:    "1.1.1990",
		Origin:       "russia",
		Profession:   "actor",
		MoviesCount:  5,
		MoviesRating: 5,
		Movies:       []models.MovieReference{{"1", "QWE", 7.5}},
		Avatar:       "qwe",
	}

	rows := pgxmock.NewRows([]string{"id", "name", "avatar"}).
		AddRow(1, actor.Movies[0].Title, actor.Avatar)
	mock.ExpectQuery("SELECT").WithArgs(actor.Name).WillReturnRows(rows)

	if _, err = movieRepo.GetFavoriteActors(actor.Name); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLikeActor(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewActorRepository(mock)
	user1 := "ilya"

	mock.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = usersRepo.LikeActor(user1, 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUnlikeActor(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewActorRepository(mock)
	user1 := "ilya"

	mock.ExpectExec("DELETE").WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err = usersRepo.UnlikeActor(user1, 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
