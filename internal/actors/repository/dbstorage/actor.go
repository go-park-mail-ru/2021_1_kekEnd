package dbstorage

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

type ActorRepository struct {
	db *pgxpool.Pool
}

func NewActorRepository(database *pgxpool.Pool) *ActorRepository {
	return &ActorRepository{
		db: database,
	}
}

func (actorStorage *ActorRepository) GetActorByID(id string) (models.Actor, error) {
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

	actor.ID = strconv.Itoa(idActor)

	moviesCount, actorMovies, err := actorStorage.getMoviesForActor(actor.Name)
	if err != nil {
		return models.Actor{}, err
	}
	actor.MoviesCount = moviesCount
	actor.Movies = actorMovies

	return actor, nil
}

func (actorStorage *ActorRepository) getMoviesForActor(name string) (int, []models.MovieReference, error) {
	var movies []models.MovieReference

	sqlStatement := `
		SELECT COUNT(*)
		FROM mdb.movie
		WHERE $1=ANY(actors)
	`

	var rowsCount int
	err := actorStorage.db.QueryRow(context.Background(), sqlStatement, name).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, movies, nil
	}
	if err != nil {
		return 0, movies, err
	}

	sqlStatement = `
		SELECT id, title, ROUND(CAST(rating AS numeric), 1)
		FROM mdb.movie
		WHERE $1=ANY(actors)
		ORDER BY rating DESC
		LIMIT $2
	`
	
	rows, err := actorStorage.db.Query(context.Background(), sqlStatement, name, _const.MoviesNumberOnActorPage)
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