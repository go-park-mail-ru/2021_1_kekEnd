package dbstorage

//import (
//	"context"
//	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
//	"github.com/pashagolub/pgxmock"
//	"testing"
//)
//
//func TestGetActorByID(t *testing.T) {
//	mock, err := pgxmock.NewConn()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer mock.Close(context.Background())
//
//	movieRepo := NewActorRepository(mock)
//	actor := &models.Actor{
//		ID: "1",
//		Name: "Ivan",
//		Biography: "from russia",
//		BirthDate: "1.1.1990",
//		Origin: "russia",
//		Profession: "actor",
//		MoviesCount: 5,
//		MoviesRating: 5,
//		Movies: []models.MovieReference{{"1", "QWE", 7.5}},
//		Avatar: "qwe",
//	}
//
//	rows := pgxmock.NewRows([]string{"id", "name", "biography", "birthdate", "origin", "profession", "avatar"}).
//	AddRow(1, actor.Name, actor.Biography, actor.BirthDate, actor.Origin, actor.Profession, actor.Avatar)
//
//	mock.ExpectQuery("SELECT").WithArgs(1, "").WillReturnRows(rows)
//
//	rows2 := pgxmock.NewRows([]string{"cnt"}).AddRow(1)
//
//	mock.ExpectQuery("SELECT").WithArgs(actor.Name).WillReturnRows(rows2)
//
//	rows3 := pgxmock.NewRows([]string{"id", "title", "rnd"}).
//	AddRow(1, actor.Movies[0].Title, actor.Movies[0].Rating)
//
//	mock.ExpectQuery("SELECT").WithArgs(actor.Name, 10).WillReturnRows(rows3)
//
//	if _, err = movieRepo.GetActorByID(actor.ID, ""); err != nil {
//		t.Errorf("error was not expected while updating stats: %s", err)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//
//	// TODO: add tests for getMoviesForActor
//}
