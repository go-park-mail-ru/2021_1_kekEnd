package localstorage

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"math"
	"strconv"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

func getHashedPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

type UserRepository struct {
	db PgxPoolIface
}

func NewUserRepository(database PgxPoolIface) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

func (storage *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := getHashedPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	sqlStatement := `
        INSERT INTO mdb.users (login, password, email)
        VALUES ($1, $2, $3)
    `

	_, errDB := storage.db.
		Exec(context.Background(), sqlStatement, user.Username, user.Password, user.Email)

	if errDB != nil {
		return errors.New("Create User Error")
	}

	return nil
}

func (storage *UserRepository) CheckEmailUnique(newEmail string) error {
	sqlStatement := `
        SELECT COUNT(*) as count
        FROM mdb.users
        WHERE email=$1
    `

	var count int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement, newEmail).
		Scan(&count)

	if err != nil || count != 0 {
		return errors.New("Email is not unique")
	}

	return nil
}

func (storage *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	sqlStatement := `
        SELECT login, password, email, img_src, movies_watched, reviews_count
        FROM mdb.users
        WHERE login=$1
    `

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, username).
		Scan(&user.Username, &user.Password,
			&user.Email, &user.Avatar,
			&user.MoviesWatched, &user.ReviewsNumber)

	if err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (storage *UserRepository) CheckPassword(password string, user *models.User) (bool, error) {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, nil
}

func (storage *UserRepository) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	if user.Username != change.Username {
		return nil, errors.New("username doesn't match")
	}

	if change.Password != "" {
		newPassword, err := getHashedPassword(change.Password)
		if err != nil {
			return nil, err
		}

		user.Password = newPassword
	}

	if change.Email != "" {
		user.Email = change.Email
	}

	if change.Avatar != "" {
		user.Avatar = change.Avatar
	}

	if change.ReviewsNumber != nil {
		user.ReviewsNumber = change.ReviewsNumber
	}

	if change.MoviesWatched != nil {
		user.MoviesWatched = change.MoviesWatched
	}

	sqlStatement := `
        UPDATE mdb.users
        SET (login, password, email, img_src, movies_watched, reviews_count) =
            ($2, $3, $4, $5, $6, $7)
        WHERE login=$1
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, user.Username,
			user.Username, user.Password,
			user.Email, user.Avatar,
			user.MoviesWatched, user.ReviewsNumber)

	if err != nil {
		return nil, errors.New("Updating user error")
	}

	return user, nil
}

func (storage *UserRepository) CheckUnsubscribed(subscriber string, user string) (error, bool) {
	sqlStatement := `
        SELECT COUNT(*) as count 
		FROM mdb.subscriptions
		WHERE user_1 = $1 AND user_2 = $2
    `

	var count int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement, subscriber, user).
		Scan(&count)

	if err != nil || count != 0 {
		return err, false
	}

	return nil, true
}

func (storage *UserRepository) Subscribe(subscriber string, user string) error {
	sqlStatement := `
        INSERT INTO mdb.subscriptions(user_1, user_2)
		VALUES ($1, $2)
    `
	_, err := storage.db.
		Exec(context.Background(), sqlStatement, subscriber, user)

	return err
}

func (storage *UserRepository) Unsubscribe(subscriber string, user string) error {
	sqlStatement := `
        DELETE FROM mdb.subscriptions(user_1, user_2)
		WHERE user_1 = $1 AND user_2 = $2
    `
	_, err := storage.db.
		Exec(context.Background(), sqlStatement, subscriber, user)

	return err
}

func (storage *UserRepository) GetModels(ids []string, limit, offset int) ([]*models.UserNoPassword, error) {
	users := make([]*models.UserNoPassword, 0)

	sqlStatement := `
        SELECT login, email, img_src, movies_watched, reviews_count
        FROM mdb.users
        WHERE login && $1 
		ORDER BY login
		LIMIT $2 OFFSET $3
    `
	rows, err := storage.db.Query(context.Background(), sqlStatement, ids, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.UserNoPassword{}
		var moviesWatched string
		var reviewsNumber string
		err = rows.Scan(&user.Username, &user.Email, &user.Avatar, &moviesWatched, &reviewsNumber)

		u, err := strconv.ParseUint(moviesWatched, 10, 64)
		if err != nil {
			return nil, err
		}
		*user.MoviesWatched = uint(u)

		u, err = strconv.ParseUint(reviewsNumber, 10, 64)
		if err != nil {
			return nil, err
		}
		*user.ReviewsNumber = uint(u)

		users = append(users, user)
	}

	return users, nil
}

func (storage *UserRepository) GetSubscribers(startIndex int, user string) (int, []*models.UserNoPassword, error) {
	subs := make([]string, 0)

	sqlStatement := `
        SELECT user_1
        FROM mdb.subscriptions
		WHERE user_2 = $1
		ORDER BY user_1
		LIMIT $2 OFFSET $3
    `

	rows, err := storage.db.Query(context.Background(), sqlStatement, user, _const.SubsPageSize, startIndex)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub string
		err = rows.Scan(&sub)
		if err != nil {
			return 0, nil, err
		}

		subs = append(subs, sub)
	}

	var rowsCount int
	sqlStatement = `
        SELECT COUNT(*)
        FROM mdb.users
        WHERE login && $1
    `
	err = storage.db.QueryRow(context.Background(), sqlStatement, subs).Scan(&rowsCount)
	if err != nil {
		return 0, nil, err
	}

	users, err := storage.GetModels(subs, _const.SubsPageSize, startIndex)

	if err != nil {
		return 0, nil, err
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.SubsPageSize))

	return pagesNumber, users, nil
}

func (storage *UserRepository) GetSubscriptions(startIndex int, user string) (int, []*models.UserNoPassword, error) {
	subs := make([]string, 0)

	sqlStatement := `
        SELECT user_2
        FROM mdb.subscriptions
		WHERE user_1 = $1
		ORDER BY user_2
		LIMIT $2 OFFSET $3
    `

	rows, err := storage.db.Query(context.Background(), sqlStatement, user, _const.SubsPageSize, startIndex)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub string
		err = rows.Scan(&sub)
		if err != nil {
			return 0, nil, err
		}

		subs = append(subs, sub)
	}

	var rowsCount int
	sqlStatement = `
        SELECT COUNT(*)
        FROM mdb.users
        WHERE login && $1
    `
	err = storage.db.QueryRow(context.Background(), sqlStatement, subs).Scan(&rowsCount)
	if err != nil {
		return 0, nil, err
	}

	users, err := storage.GetModels(subs, _const.SubsPageSize, startIndex)
	if err != nil {
		return 0, nil, err
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.SubsPageSize))

	return pagesNumber, users, nil
}