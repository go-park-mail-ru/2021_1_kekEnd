package localstorage

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

// PgxPoolIface Интерфейс для драйвера БД
type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

// RatingsRepository структура репозитория оценок
type RatingsRepository struct {
	db PgxPoolIface
}

// NewRatingsRepository инициализация репозитория оценок
func NewRatingsRepository(database PgxPoolIface) *RatingsRepository {
	return &RatingsRepository{
		db: database,
	}
}

// CreateRating создание оценки
func (storage *RatingsRepository) CreateRating(username string, movieID string, score int) error {
	sqlStatement := `
        INSERT INTO mdb.movie_rating (user_login, movie_id, rating)
        VALUES ($1, $2, $3)
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieID, score)

	if err != nil {
		return errors.New("create rating error")
	}

	return nil
}

// GetRating получение оценки
func (storage *RatingsRepository) GetRating(username string, movieID string) (models.Rating, error) {
	var rating models.Rating

	sqlStatement := `
        SELECT rating
        FROM mdb.movie_rating
        WHERE user_login = $1 AND movie_id=$2
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return models.Rating{}, err
	}

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, username, intMovieID).
		Scan(&rating.Score)

	if err != nil {
		return models.Rating{}, errors.New("Rating not found")
	}

	rating.UserID = username
	rating.MovieID = movieID

	return rating, nil
}

// DeleteRating удаление оценки
func (storage *RatingsRepository) DeleteRating(username string, movieID string) error {
	sqlStatement := `
        DELETE FROM mdb.movie_rating
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieID)

	if err != nil {
		return errors.New("delete rating error")
	}

	return nil
}

// UpdateRating обновление оценки
func (storage *RatingsRepository) UpdateRating(username string, movieID string, score int) error {
	sqlStatement := `
        UPDATE mdb.movie_rating
        SET rating = $3
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieID, score)

	if err != nil {
		return errors.New("update rating error")
	}

	return nil
}

// GetFeed получение новостей об оценках
func (storage *RatingsRepository) GetFeed(users []models.UserNoPassword) ([]models.RatingFeedItem, error) {
	feed := make([]models.RatingFeedItem, 0)

	subs := make([]string, len(users))
	for _, u := range users {
		subs = append(subs, u.Username)
	}

	sqlStatement := `
		SELECT rtg.user_login, us.img_src, rtg.movie_id, mv.title, rtg.rating, rtg.creation_date
		FROM mdb.movie_rating rtg
		JOIN mdb.users us ON rtg.user_login=us.login
		JOIN mdb.movie mv ON rtg.movie_id=mv.id
		WHERE rtg.user_login=ANY($1) AND rtg.creation_date >= NOW() - INTERVAL '48 HOURS'
		ORDER BY rtg.creation_date DESC
	`

	rows, err := storage.db.Query(context.Background(), sqlStatement, subs)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		feedItem := models.RatingFeedItem{
			ItemType: "rating",
		}
		rating := models.Rating{}
		var movieID int
		err := rows.Scan(&feedItem.Username, &feedItem.Avatar, &movieID, &feedItem.MovieTitle, &rating.Score, &feedItem.Date)
		if err != nil {
			return nil, err
		}
		rating.MovieID = strconv.Itoa(movieID)
		feedItem.Rating = rating
		feed = append(feed, feedItem)
	}

	return feed, nil
}
