package dbstorage

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"strconv"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
}

type ActorRepository struct {
	db PgxPoolIface
}

func NewActorRepository(database PgxPoolIface) *ActorRepository {
	return &ActorRepository{
		db: database,
	}
}

func (actorStorage *ActorRepository) GetActorByID(id string, username string) (models.Actor, error) {
	var actor models.Actor

	sqlStatement := `
        SELECT id, name, biography, birthdate, origin, profession, avatar
        FROM mdb.actors
        WHERE id=$1
    `

	idActor, err := strconv.Atoi(id)
	if err != nil {
		return actor, err
	}

	err = actorStorage.db.
		QueryRow(context.Background(), sqlStatement, idActor).Scan(&idActor, &actor.Name, &actor.Biography,
			&actor.BirthDate, &actor.Origin, &actor.Profession, &actor.Avatar)

	if err != nil {
		return actor, err
	}

	sqlStatementLiked := `
		SELECT COUNT(*) FROM mdb.favorite_actors WHERE user_login=$1 AND actor_id=$2
	`
	var rowsCount int
	err = actorStorage.db.QueryRow(context.Background(), sqlStatementLiked, username, idActor).Scan(&rowsCount)
	isLiked := true
	if err != nil {
		return models.Actor{}, err
	}
	if rowsCount == 0 {
		isLiked = false
	}
	actor.IsLiked = isLiked

	actor.ID = strconv.Itoa(idActor)

	moviesCount, actorMovies, err := actorStorage.getMoviesForActor(actor.ID)
	if err != nil {
		return models.Actor{}, err
	}
	actor.MoviesCount = moviesCount
	actor.Movies = actorMovies

	return actor, nil
}

func (actorStorage *ActorRepository) getMoviesForActor(id string) (int, []models.MovieReference, error) {
	var movies []models.MovieReference

	sqlStatement := `
		SELECT COUNT(*) as cnt
		FROM mdb.movie_actors
		WHERE actor_id=$1
	`

	var rowsCount int
	err := actorStorage.db.QueryRow(context.Background(), sqlStatement, id).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, movies, nil
	}
	if err != nil {
		return 0, movies, err
	}

	sqlStatement = `
		SELECT id, title, ROUND(CAST(rating AS numeric), 1) as rnd
		FROM mdb.movie mv
		JOIN movie_actors mvac ON mv.id = mvac.movie_id AND mvac.actor_id=$1
		GROUP BY mv.id
		ORDER BY rating DESC
		LIMIT $2
	`

	rows, err := actorStorage.db.Query(context.Background(), sqlStatement, id, _const.MoviesNumberOnActorPage)
	if err != nil {
		return 0, movies, err
	}
	defer rows.Close()

	for rows.Next() {
		movieReference := &models.MovieReference{}
		var id int
		err = rows.Scan(&id, &movieReference.Title, &movieReference.Rating)
		if err != nil && err != sql.ErrNoRows {
			return 0, []models.MovieReference{}, err
		}
		movieReference.ID = strconv.Itoa(id)
		movies = append(movies, *movieReference)
	}

	return rowsCount, movies, nil
}

func (actorStorage *ActorRepository) CreateActor(actor models.Actor) error {
	return nil
}

func (actorStorage *ActorRepository) EditActor(actor models.Actor) (models.Actor, error) {
	return models.Actor{}, nil
}

func (actorStorage *ActorRepository) LikeActor(username string, actorID int) error {
	sqlStatement := `
		INSERT INTO mdb.favorite_actors VALUES($1, $2)
	`

	_, err := actorStorage.db.Exec(context.Background(), sqlStatement, username, actorID)
	if err != nil {
		return err
	}
	return nil
}

func (actorStorage *ActorRepository) UnlikeActor(username string, actorID int) error {
	sqlStatement := `
		DELETE FROM mdb.favorite_actors WHERE user_login=$1 AND actor_id=$2
	`

	_, err := actorStorage.db.Exec(context.Background(), sqlStatement, username, actorID)
	if err != nil {
		return err
	}
	return nil
}
