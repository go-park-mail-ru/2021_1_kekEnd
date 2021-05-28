package localstorage

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
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

// ReviewRepository структура репозитория рецензий
type ReviewRepository struct {
	db PgxPoolIface
}

// NewReviewRepository инициализация репозитория рецензий
func NewReviewRepository(database PgxPoolIface) *ReviewRepository {
	return &ReviewRepository{
		db: database,
	}
}

func convertReviewTypeFromIntToStr(reviewType int) models.ReviewType {
	switch reviewType {
	case 1:
		return "positive"
	case 0:
		return "neutral"
	case -1:
		return "negative"
	default:
		return ""
	}
}

func convertReviewTypeFromStrToInt(reviewType models.ReviewType) int {
	switch reviewType {
	case "positive":
		return 1
	case "neutral":
		return 0
	case "negative":
		return -1
	default:
		return -100
	}
}

// CreateReview создание рецензии
func (storage *ReviewRepository) CreateReview(review *models.Review) error {
	sqlStatement := `
        INSERT INTO mdb.users_review (user_login, movie_id, review_type, title, content)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "id";
    `
	var newID int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement,
			review.Author, review.MovieID,
			convertReviewTypeFromStrToInt(review.ReviewType), review.Title,
			review.Content).
		Scan(&newID)

	if err != nil {
		return errors.New("create review error")
	}

	review.ID = strconv.Itoa(newID)

	return nil
}

// GetUserReviews полученить рецензии пользователя
func (storage *ReviewRepository) GetUserReviews(username string) ([]*models.Review, error) {
	var reviews []*models.Review

	sqlStatement := `
        SELECT ur.id, ur.movie_id, mv.title, ur.review_type, ur.title, ur.content
        FROM mdb.users_review ur JOIN mdb.movie mv ON ur.movie_id = mv.id
        WHERE ur.user_login = $1
    `

	rows, err := storage.db.
		Query(context.Background(), sqlStatement, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := &models.Review{}
		var newID int
		var newMovieID int
		var newReviewType int
		err = rows.Scan(&newID, &newMovieID, &review.MovieTitle, &newReviewType, &review.Title, &review.Content)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		review.ID = strconv.Itoa(newID)
		review.Author = username
		review.MovieID = strconv.Itoa(newMovieID)
		review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// GetMovieReviews получить рецензии фильма
func (storage *ReviewRepository) GetMovieReviews(movieID string, startInd int) (int, []*models.Review, error) {
	var reviews []*models.Review

	sqlStatement := `
        SELECT COUNT(*)
        FROM mdb.users_review
		WHERE movie_id = $1;
    `
	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return 0, nil, err
	}

	var rowsCount int
	err = storage.db.QueryRow(context.Background(), sqlStatement, intMovieID).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, reviews, nil
	}
	if err != nil {
		return 0, nil, err
	}

	sqlStatement = `
        SELECT id, user_login, review_type, title, content
        FROM mdb.users_review
        WHERE movie_id=$1
        ORDER BY creation_date
        LIMIT $2 OFFSET $3
    `

	rows, err := storage.db.Query(context.Background(), sqlStatement, intMovieID, constants.ReviewsPageSize, startInd)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := &models.Review{}
		var newID int
		var newReviewType int
		err = rows.Scan(&newID, &review.Author, &newReviewType, &review.Title, &review.Content)
		if err != nil && err != sql.ErrNoRows {
			return 0, nil, err
		}

		review.ID = strconv.Itoa(newID)
		review.MovieID = movieID
		review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

		reviews = append(reviews, review)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / constants.ReviewsPageSize))

	return pagesNumber, reviews, nil
}

// GetUserReviewForMovie получить рецензию пользователя к фильму
func (storage *ReviewRepository) GetUserReviewForMovie(username string, movieID string) (*models.Review, error) {
	var review models.Review

	sqlStatement := `
        SELECT id, review_type, title, content
        FROM mdb.users_review
        WHERE user_login = $1 AND movie_id=$2
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return nil, err
	}

	var newID int
	var newReviewType int

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, username, intMovieID).
		Scan(&newID, &newReviewType, &review.Title, &review.Content)

	if err != nil {
		return nil, errors.New("review not found")
	}

	review.ID = strconv.Itoa(newID)
	review.Author = username
	review.MovieID = movieID
	review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

	return &review, nil
}

// EditUserReviewForMovie изменить рецензию пользователя
func (storage *ReviewRepository) EditUserReviewForMovie(review *models.Review) error {
	sqlStatement := `
        UPDATE mdb.users_review
        SET (review_type, title, content) =
            ($3, $4, $5)
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieID, err := strconv.Atoi(review.MovieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, review.Author, intMovieID,
			convertReviewTypeFromStrToInt(review.ReviewType), review.Title,
			review.Content)

	if err != nil {
		return errors.New("update review error")
	}

	return nil
}

// DeleteUserReviewForMovie удалить рецензнию пользователя
func (storage *ReviewRepository) DeleteUserReviewForMovie(username string, movieID string) error {
	_, err := storage.GetUserReviewForMovie(username, movieID)
	if err != nil {
		return err
	}

	sqlStatement := `
        DELETE FROM mdb.users_review
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieID)

	if err != nil {
		return errors.New("delete review error")
	}

	return nil
}

// GetFeed получие новостей о рецензиях
func (storage *ReviewRepository) GetFeed(users []models.UserNoPassword) ([]models.ReviewFeedItem, error) {
	feed := make([]models.ReviewFeedItem, 0)

	subs := make([]string, len(users))
	for _, u := range users {
		subs = append(subs, u.Username)
	}

	sqlStatement := `
        SELECT rws.user_login, us.img_src, rws.title, rws.content, rws.review_type, rws.creation_date
        FROM mdb.users_review rws
		JOIN mdb.users us ON rws.user_login=us.login
        WHERE rws.user_login=ANY($1) AND rws.creation_date >= NOW() - INTERVAL '48 HOURS'
        ORDER BY rws.creation_date DESC
    `

	rows, err := storage.db.Query(context.Background(), sqlStatement, subs)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		feedItem := models.ReviewFeedItem{
			ItemType: "review",
		}
		review := models.Review{}
		var reviewTypeInt int
		err := rows.Scan(&feedItem.Username, &feedItem.Avatar, &review.Title, &review.Content, &reviewTypeInt, &feedItem.Date)
		if err != nil {
			return nil, err
		}
		review.ReviewType = convertReviewTypeFromIntToStr(reviewTypeInt)
		feedItem.Review = review
		feed = append(feed, feedItem)
	}

	return feed, nil
}
