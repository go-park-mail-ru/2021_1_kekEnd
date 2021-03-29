package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"strconv"
    "github.com/jackc/pgx/v4/pgxpool"
    "context"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"

)

type ReviewRepository struct {
    db *pgxpool.Pool
}

func NewReviewRepository(database *pgxpool.Pool) *ReviewRepository {
    return &ReviewRepository{
        db: database,
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
                    review.ReviewType, review.Title,
                    review.Content).
           Scan(&newID)

    if err != nil {
        return errors.New("Create Review Error")
    }

    review.ID = strconv.Itoa(newID)

    return nil
}

func (storage *ReviewRepository) GetUserReviews(username string) []*models.Review {
    var reviews []*models.Review

    sqlStatement := `
        SELECT id, movie_id, review_type, title, content
        FROM mdb.users_review
        WHERE u.login = $1
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
        err = rows.Scan(&newID, &newMovieID, &review.ReviewType, &review.Title, &review.Content)
        if err != nil {
            return nil
        }

        review.ID = strconv.Itoa(newID)
        review.Author = username

        reviews = append(reviews, review)
    }

    return reviews
}

func (storage *ReviewRepository) GetMovieReviews(movieID string, page int) (int, []*models.Review) {
    var reviews []*models.Review

    sqlStatement := `
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

	startIndex := (page - 1) * _const.ReviewsPageSize
	endIndex := startIndex + _const.ReviewsPageSize

    rows, err := storage.db.
           Query(context.Background(), sqlStatement, intMovieId, endIndex - startIndex, startIndex)
    if err != nil {
        return 0, nil
    }
	defer rows.Close()

    for rows.Next() {
        review := &models.Review{}
        var newID int
        err = rows.Scan(&newID, &review.Author, &review.ReviewType, &review.Title, &review.Content)
        if err != nil {
            return 0, nil
        }

        review.ID = strconv.Itoa(newID)
        review.MovieID = movieID

        reviews = append(reviews, review)
    }

    if len(reviews) == 0 {
    	return 0, nil
    }

    return len(reviews), reviews
}

func (storage *ReviewRepository) GetUserReviewForMovie(username string, movieID string) (*models.Review, error)  {
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

    err = storage.db.
           QueryRow(context.Background(), sqlStatement, username, intMovieId).
           Scan(&newID, &review.ReviewType, &review.Title, &review.Content)

    if err != nil {
        return nil, errors.New("Review not found")
    }

    review.Author = username
    review.MovieID = movieID

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
           review.ReviewType, review.Title,
           review.Content)

    if err != nil {
        return errors.New("Update Review Error")
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
        return errors.New("Delete Review Error")
    }

    return nil
}
