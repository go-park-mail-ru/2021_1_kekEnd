package dbstorage

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
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

	return actor, nil
}
