package localstorage

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgconn"
	"strconv"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

type RatingsRepository struct {
	db PgxPoolIface
}

func NewRatingsRepository(database PgxPoolIface) *RatingsRepository {
	return &RatingsRepository{
		db: database,
	}
}

func (storage *RatingsRepository) CreateRating(username string, movieID string, score int) error {
	sqlStatement := `
        INSERT INTO mdb.movie_rating (user_login, movie_id, rating)
        VALUES ($1, $2, $3)
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieId, score)

	if err != nil {
		return errors.New("create rating error")
	}

	return nil
}

func (storage *RatingsRepository) GetRating(username string, movieID string) (models.Rating, error) {
	var rating models.Rating

	sqlStatement := `
        SELECT rating
        FROM mdb.movie_rating
        WHERE user_login = $1 AND movie_id=$2
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return models.Rating{}, err
	}

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, username, intMovieId).
		Scan(&rating.Score)

	if err != nil {
		return models.Rating{}, errors.New("Rating not found")
	}

	rating.UserID = username
	rating.MovieID = movieID

	return rating, nil
}

func (storage *RatingsRepository) DeleteRating(username string, movieID string) error {
	sqlStatement := `
        DELETE FROM mdb.movie_rating
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieId)

	if err != nil {
		return errors.New("delete rating error")
	}

	return nil
}

func (storage *RatingsRepository) UpdateRating(username string, movieID string, score int) error {
	sqlStatement := `
        UPDATE mdb.movie_rating
        SET rating = $3
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieId, score)

	if err != nil {
		return errors.New("update rating error")
	}

	return nil
}
