package dbstorage

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func TestCreatePlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)
	user1 := "ilya"
	playlistName := "newPlaylist"

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery("INSERT").WithArgs(playlistName, user1, false).WillReturnRows(rows)

	mock.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = usersRepo.CreatePlaylist(user1, playlistName, false); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"id", "name", "movies"}).
		AddRow(1, "newPlaylist", map[int]string{})
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	if _, err = usersRepo.GetPlaylist(1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlaylists(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"id", "name", "movies"}).
		AddRow(1, "newPlaylist", map[int]string{})
	mock.ExpectQuery("SELECT").WithArgs("ilya").WillReturnRows(rows)

	if _, err = usersRepo.GetPlaylists("ilya"); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlaylistsInfo(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"id", "name", "movie_id"}).
		AddRow(1, "newPlaylist", 1)
	mock.ExpectQuery("SELECT").WithArgs(1, "ilya").WillReturnRows(rows)

	if _, err = usersRepo.GetPlaylistsInfo("ilya", 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCanUserUpdatePlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"username"}).
		AddRow("ilya")
	mock.ExpectQuery("SELECT").WithArgs(1, "ilya").WillReturnRows(rows)

	if err = usersRepo.CanUserUpdatePlaylist("ilya", 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteAllUserFromPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"username"}).
		AddRow("ilya")
	mock.ExpectQuery("DELETE").WithArgs(1, "ilya").WillReturnRows(rows)

	rows2 := pgxmock.NewRows([]string{"username"}).
		AddRow("ilya")
	mock.ExpectQuery("DELETE").WithArgs(1, "ilya").WillReturnRows(rows2)

	if err = usersRepo.DeleteAllUserFromPlaylist("ilya", 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"isShared"}).
		AddRow(false)
	mock.ExpectQuery("UPDATE").WithArgs(1, "plst1", false).WillReturnRows(rows)

	if err = usersRepo.UpdatePlaylist("ilya", 1, "plst1", false); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeletePlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	mock.ExpectExec("DELETE").WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err = usersRepo.DeletePlaylist(1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCanUserUpdateMovieInPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"username"}).
		AddRow("ilya")
	mock.ExpectQuery("SELECT").WithArgs(1, "ilya").WillReturnRows(rows)

	if err = usersRepo.CanUserUpdateMovieInPlaylist("ilya", 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddMovieToPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	mock.ExpectExec("INSERT").WithArgs(1, 1, "ilya").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = usersRepo.AddMovieToPlaylist("ilya", 1, 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteMovieFromPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	mock.ExpectExec("DELETE").WithArgs(1, 1).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err = usersRepo.DeleteMovieFromPlaylist("ilya", 1, 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCanUserUpdateUsersInPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	rows := pgxmock.NewRows([]string{"username"}).
		AddRow("ilya")
	mock.ExpectQuery("SELECT").WithArgs(1, "ilya").WillReturnRows(rows)

	if err = usersRepo.CanUserUpdateUsersInPlaylist("ilya", 1); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddUserToPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	mock.ExpectExec("INSERT").WithArgs("ivan", 1).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err = usersRepo.AddUserToPlaylist("ilya", 1, "ivan"); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUserFromPlaylist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	usersRepo := NewPlaylistsRepository(mock)

	mock.ExpectExec("DELETE").WithArgs("ivan", 1).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err = usersRepo.DeleteUserFromPlaylist("ilya", 1, "ivan"); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
