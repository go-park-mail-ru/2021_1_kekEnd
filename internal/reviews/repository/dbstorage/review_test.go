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

	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	movieRepo := NewReviewRepository(mock)

	mock.ExpectExec("INSERT INTO").WithArgs(review.Author, review.MovieID, review.ReviewType, review.Title, review.Content).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = movieRepo.CreateReview(review); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserReviews(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewReviewRepository(mock)
	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	rows := pgxmock.NewRows([]string{"id", "movie_id", "review_type", "title", "content"}).
		AddRow(1, 1, 0, review.Title, review.Content)

	mock.ExpectQuery("SELECT").WithArgs(review.Author).WillReturnRows(rows)

	// now we execute our method
	if _, err = movieRepo.GetUserReviews(review.Author); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMovieReviews(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewReviewRepository(mock)
	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	rows1 := pgxmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery("SELECT").WillReturnRows(rows1)

	rows2 := pgxmock.NewRows([]string{"id", "movie_id", "review_type", "title", "content"}).
		AddRow(1, 1, 0, review.Title, review.Content)
	mock.ExpectQuery("SELECT").WithArgs(1, 3, 1).WillReturnRows(rows2)

	// now we execute our method
	if _, _, err = movieRepo.GetMovieReviews(review.MovieID, 1); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserReviewForMovie(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewReviewRepository(mock)
	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	rows2 := pgxmock.NewRows([]string{"id", "review_type", "title", "content"}).
		AddRow(1, 0, review.Title, review.Content)
	mock.ExpectQuery("SELECT").WithArgs(review.Author, 1).WillReturnRows(rows2)

	// now we execute our method
	if _, err = movieRepo.GetUserReviewForMovie(review.Author, review.MovieID); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEditUserReviewForMovie(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewReviewRepository(mock)
	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	mock.ExpectExec("UPDATE").WithArgs(review.Author, 1, 0, review.Title, review.Content).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// now we execute our method
	if err = movieRepo.EditUserReviewForMovie(review); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUserReviewForMovie(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	movieRepo := NewReviewRepository(mock)
	review := &models.Review{
		ID:         "1",
		Title:      "Goog",
		Content:    "good film",
		ReviewType: "neutral",
		Author:     "ILYA",
		MovieID:    "1",
	}

	mock.ExpectExec("DELETE").WithArgs(review.Author, 1).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// now we execute our method
	if err = movieRepo.DeleteUserReviewForMovie(review.Author, review.MovieID); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
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

	ratingRepo := NewReviewRepository(mock)

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
