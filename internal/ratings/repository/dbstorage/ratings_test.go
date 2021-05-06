package localstorage

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/pashagolub/pgxmock"
)

func TestCreateRating(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	username := "ilya"
	movieID := "1"
	score := 10

	movieRepo := NewRatingsRepository(mock)

	mock.ExpectExec("INSERT INTO").WithArgs(username, 1, score).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = movieRepo.CreateRating(username, movieID, score); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRating(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	username := "ilya"
	movieID := "1"

	movieRepo := NewRatingsRepository(mock)

	rows := pgxmock.NewRows([]string{"rating"}).AddRow(5)

	mock.ExpectQuery("SELECT").WithArgs(username, 1).WillReturnRows(rows)

	if _, err = movieRepo.GetRating(username, movieID); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteRating(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	username := "ilya"
	movieID := "1"

	movieRepo := NewRatingsRepository(mock)

	mock.ExpectExec("DELETE").WithArgs(username, 1).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err = movieRepo.DeleteRating(username, movieID); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateRating(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	username := "ilya"
	movieID := "1"
	score := 2

	movieRepo := NewRatingsRepository(mock)

	mock.ExpectExec("UPDATE mdb.movie_rating").WithArgs(username, 1, score).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err = movieRepo.UpdateRating(username, movieID, score); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetFeed(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	user := []models.UserNoPassword{
		models.UserNoPassword{
			Username: "ilya",
		},
	}

	ratingRepo := NewRatingsRepository(mock)

	rows := pgxmock.NewRows([]string{"user_login", "img_src", "movie_id", "title", "rating", "creation_date"}).
		AddRow("ilya", "", "", "", "", "")

	mock.ExpectQuery("SELECT").WithArgs("ilya").WillReturnRows(rows)

	if _, err = ratingRepo.GetFeed(user); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
