package localstorage

import (
    "errors"
    "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
    "github.com/jackc/pgx/v4/pgxpool"
    "strconv"
    "context"
    "fmt"
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
        INSERT INTO mdb.users_review (user_id, movie_id, review_type, title, content)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING "id";
    `

    var myid int
    errDB := storage.db.
           QueryRow(context.Background(), sqlStatement,
                    1, 2,
                    review.ReviewType, review.Title,
                    review.Content).
           Scan(&myid)

    if errDB != nil {
        return errors.New("Create Review Error")
    }

    return nil
}


// Надо переделать, чтобы передавался не ник, а ID пользователя
func (storage *ReviewRepository) GetUserReviews(username string) []*models.Review {
    var reviews []*models.Review

    sqlStatement := `
        SELECT user_id, movie_id, review_type, title, content
        FROM mdb.users_review ur JOIN mdb.users u ON ur.user_id = u.user_id
        AND u.login = $1
    `

    rows, err := storage.db.
           Query(context.Background(), sqlStatement, username)

    if err != nil {
        return nil
    }

    for rows.Next() {
        review := &models.Review{}
        var newUserID int
        var newMovieID int
        err = rows.Scan(&newUserID, &newMovieID, &review.ReviewType, &review.Title, &review.Content)
        if err != nil {
            return nil
        }

        // review.UserID = strconv.Itoa(newUserID)
        review.MovieID = strconv.Itoa(newMovieID)

        reviews = append(reviews, review)
    }

    return reviews
}

func (storage *ReviewRepository) GetMovieReviews(movieID string) []*models.Review {
    var reviews []*models.Review

    sqlStatement := `
        SELECT user_id, movie_id, review_type, title, content
        FROM mdb.users_review
        WHERE movie_id=$1
    `

    intMovieId, err := strconv.Atoi(movieID)
    if err != nil {
        return nil
    }

    rows, err := storage.db.
           Query(context.Background(), sqlStatement, intMovieId)

    if err != nil {
        return nil
    }

    for rows.Next() {
        review := &models.Review{}
        var newUserID int
        var newMovieID int
        err = rows.Scan(&newUserID, &newMovieID, &review.ReviewType, &review.Title, &review.Content)
        if err != nil {
            return nil
        }

        // review.userID = strconv.Itoa(userID)
        review.MovieID = strconv.Itoa(newMovieID)

        reviews = append(reviews, review)
    }

    return reviews
}

func (storage *ReviewRepository) GetUserReviewForMovie(username string, movieID string) ([]*models.Review, error)  {
    var reviews []*models.Review

    sqlStatement := `
        SELECT user_id, movie_id, review_type, title, content
        FROM mdb.users_review ur JOIN mdb.users u ON ur.user_id = u.user_id
        WHERE u.login = $1 AND ur.movie_id=$2
    `

    intMovieId, err := strconv.Atoi(movieID)
    if err != nil {
        return nil, err
    }

    rows, err := storage.db.
           Query(context.Background(), sqlStatement, username, intMovieId)

    if err != nil {
        return nil, errors.New("Review not found")
    }

    for rows.Next() {
        review := &models.Review{}
        var newUserID int
        var newMovieID int
        err = rows.Scan(&newUserID, &newMovieID, &review.ReviewType, &review.Title, &review.Content)
        if err != nil {
            return nil, err
        }

        // review.userID = strconv.Itoa(userID)
        review.MovieID = strconv.Itoa(newMovieID)

        reviews = append(reviews, review)
    }

    return reviews, nil
}

func (storage *ReviewRepository) EditUserReviewForMovie(review *models.Review) error {
    return nil
}

func (storage *ReviewRepository) DeleteUserReviewForMovie(username string, movieID string) error {
    sqlStatement := `
        DELETE
        FROM mdb.users_review ur JOIN mdb.users u ON ur.user_id = u.user_id
        WHERE u.login = $1 AND ur.movie_id=$2
        RETURNING "id";
    `

    intMovieId, err := strconv.Atoi(movieID)
    if err != nil {
        return err
    }

    var newID int
    errDB := storage.db.
           QueryRow(context.Background(), sqlStatement, username, intMovieId).
           Scan(&newID)
    fmt.Println(errDB)

    if errDB != nil {
        return errors.New("Create Review Error")
    }

    return nil
}
