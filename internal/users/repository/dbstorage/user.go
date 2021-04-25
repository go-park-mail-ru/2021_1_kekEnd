package localstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	// "fmt"
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

func (storage *UserRepository) Subscribe(subscriber string, user string) error {
	sqlStatement := `
        SELECT COUNT(*) as count 
		FROM mdb.subscriptions
		WHERE user_1 = $1 AND user_2 = $2
    `

	var count int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement, subscriber, user).
		Scan(&count)

	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("%s is already subscribed to %s", subscriber, user)
	}

	sqlStatement = `
        INSERT INTO mdb.subscriptions(user_1, user_2)
		VALUES ($1, $2)
    `
	_, err = storage.db.
		Exec(context.Background(), sqlStatement, subscriber, user)

	return err
}

func (storage *UserRepository) Unsubscribe(subscriber *models.User, user *models.User) error {
	return nil
}

func (storage *UserRepository) GetSubscribers(user *models.User) []models.User {
	return nil
}

func (storage *UserRepository) GetSubscriptions(user *models.User) []models.User {
	return nil
}
