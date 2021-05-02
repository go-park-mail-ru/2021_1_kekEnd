package localstorage

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	// "fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)
	user := &models.User{
		Username:      "login",
		Email:         "email",
		Password:      "password",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WithArgs(user.Username, user.Password, user.Email).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = usersRepo.CreateUser(user); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCheckEmailUnique(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)
	email := "ilya228@mail.ru"

	rows := pgxmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery("SELECT").WithArgs(email).WillReturnRows(rows)

	// now we execute our method
	if err = usersRepo.CheckEmailUnique(email); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)
	user := &models.User{
		Username:      "login",
		Email:         "email",
		Password:      "password",
		Avatar:        _const.DefaultAvatarPath,
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	rows := pgxmock.NewRows([]string{"login", "password", "email", "img_src", "movies_watched", "reviews_count"}).
	AddRow(user.Username, user.Password, user.Email, user.Avatar, user.MoviesWatched, user.ReviewsNumber)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(user.Username).WillReturnRows(rows)
	mock.ExpectCommit()

	// now we execute our method
	if _, err = usersRepo.GetUserByUsername(user.Username); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)
	user1 := &models.User{
		Username:      "login",
		Email:         "email1",
		Password:      "password",
		Avatar:        _const.DefaultAvatarPath,
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	user2 := models.User{
		Username:      "login",
		Email:         "email2",
		Password:      "password",
		Avatar:        _const.DefaultAvatarPath,
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	// now we execute our method
	if _, err = usersRepo.UpdateUser(user1, user2); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCheckUnsubscribed(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())


	usersRepo := NewUserRepository(mock)
	user1 := "lol"
	user2 := "kek"

	mock.ExpectBegin()
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 0))
	mock.ExpectCommit()

	// now we execute our method
	if err, _ = usersRepo.CheckUnsubscribed(user1, user2); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


func TestSubscribe(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())


	usersRepo := NewUserRepository(mock)
	user1 := "lol"
	user2 := "kek"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = usersRepo.Subscribe(user1, user2); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUnsubscribe(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())


	usersRepo := NewUserRepository(mock)
	user1 := "lol"
	user2 := "kek"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = usersRepo.Unsubscribe(user1, user2); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetModels(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)

	subs := []string{"lol", "kek"}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 2))
	mock.ExpectCommit()

	// now we execute our method
	if _, err = usersRepo.GetModels(subs, 0); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSubscribers(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)

	user := "lol"

	mock.ExpectBegin()
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 0))
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 0))
	mock.ExpectCommit()

	// now we execute our method
	if _, subs, err := usersRepo.GetSubscribers(0, user); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
		assert.Equal(t, subs, []*models.UserNoPassword{})
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSubscriptions(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewUserRepository(mock)

	user := "lol"

	mock.ExpectBegin()
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 0))
	mock.ExpectExec("SELECT").WillReturnResult(pgxmock.NewResult("SELECT", 0))
	mock.ExpectCommit()

	// now we execute our method
	if _, subs, err := usersRepo.GetSubscriptions(0, user); err == nil {
		t.Errorf("error was not expected while updating stats: %s", err)
		assert.Equal(t, subs, []*models.UserNoPassword{})
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}