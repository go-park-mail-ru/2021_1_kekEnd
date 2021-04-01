package localstorage

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"strconv"
)

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(database *pgxpool.Pool) *ReviewRepository {
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

func (storage *ReviewRepository) GetUserReviews(username string) []*models.Review {
	var reviews []*models.Review

	sqlStatement := `
        SELECT id, movie_id, review_type, title, content
        FROM mdb.users_review
        WHERE user_login = $1
    `

	rows, err := storage.db.
		Query(context.Background(), sqlStatement, username)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		review := &models.Review{}
		var newID int
		var newMovieID int
		var newReviewType int
		err = rows.Scan(&newID, &newMovieID, &newReviewType, &review.Title, &review.Content)
		if err != nil {
			return nil
		}

		review.ID = strconv.Itoa(newID)
		review.Author = username
		review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

		reviews = append(reviews, review)
	}

	return reviews
}

func (storage *ReviewRepository) GetMovieReviews(movieID string, startInd int) (int, []*models.Review) {
	var reviews []*models.Review

	sqlStatement := `
        SELECT COUNT(*)
        FROM mdb.users_review
    `

	var rowsCount int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement).Scan(&rowsCount)
	if err != nil {
		return 0, nil
	}

	sqlStatement = `
        SELECT id, user_login, review_type, title, content
        FROM mdb.users_review
        WHERE movie_id=$1
        ORDER BY creation_date
        LIMIT $2 OFFSET $3
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return 0, nil
	}

	rows, err := storage.db.
		Query(context.Background(), sqlStatement, intMovieId, _const.ReviewsPageSize, startInd)
	if err != nil {
		return 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		review := &models.Review{}
		var newID int
		var newReviewType int
		err = rows.Scan(&newID, &review.Author, &newReviewType, &review.Title, &review.Content)
		if err != nil {
			return 0, nil
		}

		review.ID = strconv.Itoa(newID)
		review.MovieID = movieID
		review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

		reviews = append(reviews, review)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.ReviewsPageSize))

	return pagesNumber, reviews
}

func (storage *ReviewRepository) GetUserReviewForMovie(username string, movieID string) (*models.Review, error) {
	var review models.Review

	sqlStatement := `
        SELECT id, review_type, title, content
        FROM mdb.users_review
        WHERE user_login = $1 AND movie_id=$2
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return nil, err
	}

	var newID int
	var newReviewType int

	err = storage.db.
		QueryRow(context.Background(), sqlStatement, username, intMovieId).
		Scan(&newID, &newReviewType, &review.Title, &review.Content)

	if err != nil {
		return nil, errors.New("review not found")
	}

	review.Author = username
	review.MovieID = movieID
	review.ReviewType = convertReviewTypeFromIntToStr(newReviewType)

	return &review, nil
}

func (storage *ReviewRepository) EditUserReviewForMovie(review *models.Review) error {
	sqlStatement := `
        UPDATE mdb.users_review
        SET (review_type, title, content) =
            ($3, $4, $5)
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieId, err := strconv.Atoi(review.MovieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, review.Author, intMovieId,
			convertReviewTypeFromStrToInt(review.ReviewType), review.Title,
			review.Content)

	if err != nil {
		return errors.New("update review error")
	}

	return nil
}

func (storage *ReviewRepository) DeleteUserReviewForMovie(username string, movieID string) error {
	_, err := storage.GetUserReviewForMovie(username, movieID)
	if err != nil {
		return err
	}

	sqlStatement := `
        DELETE FROM mdb.users_review
        WHERE user_login=$1 AND movie_id=$2;
    `

	intMovieId, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	_, err = storage.db.
		Exec(context.Background(), sqlStatement, username, intMovieId)

	if err != nil {
		return errors.New("delete review error")
	}

	return nil
}