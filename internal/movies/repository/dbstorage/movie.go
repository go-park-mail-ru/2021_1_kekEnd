package localstorage

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	pgx "github.com/jackc/pgx/v4"
	"math"
	"strconv"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close(context.Context) error
}

type MovieRepository struct {
	db PgxIface
}

func NewMovieRepository(database PgxIface) *MovieRepository {
	return &MovieRepository{
		db: database,
	}
}

func (storage *MovieRepository) CreateMovie(movie *models.Movie) error {
	return nil
}

func (storage *MovieRepository) GetMovieByID(id string) (*models.Movie, error) {
	tx, err := storage.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()

	var movie models.Movie

	sqlStatement := `
        SELECT id, title, description, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1) AS rating, rating_count
        FROM mdb.movie
        WHERE id=$1
    `

	idFilm, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var actorsNames []string
	err = tx.
		QueryRow(context.Background(), sqlStatement, idFilm).Scan(&idFilm,
		&movie.Title, &movie.Description,
		&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
		&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
		&movie.Montage, &movie.Budget, &movie.Duration, &actorsNames, &movie.Poster,
		&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)

	if err != nil {
		return nil, err
	}

	movie.Actors, err = storage.getActorsData(actorsNames)
	if err != nil {
		return nil, err
	}

	movie.ID = strconv.Itoa(idFilm)

	return &movie, nil
}

func (storage *MovieRepository) GetBestMovies(startIndex int) (int, []*models.Movie, error) {
	tx, err := storage.db.Begin(context.Background())
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()

	var bestMovies []*models.Movie

	sqlStatement := `
		SELECT movies_count
		FROM mdb.meta
		ORDER BY version DESC
	`

	var rowsCount int
	err = tx.QueryRow(context.Background(), sqlStatement).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, bestMovies, nil
	}
	if err != nil {
		return 0, nil, err
	}

	if rowsCount > _const.MoviesTop100Size {
		rowsCount = _const.MoviesTop100Size
	}

	sqlStatement = `
        SELECT id, title, description, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie
        ORDER BY rating DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := tx.Query(context.Background(), sqlStatement, _const.MoviesPageSize, startIndex)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		movie := &models.Movie{}
		var id int
		var actorsNames []string
		err = rows.Scan(&id, &movie.Title, &movie.Description,
			&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
			&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
			&movie.Montage, &movie.Budget, &movie.Duration, &actorsNames, &movie.Poster,
			&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)
		if err != nil && err != sql.ErrNoRows {
			return 0, nil, err
		}

		movie.ID = strconv.Itoa(id)
		movie.Actors, err = storage.getActorsData(actorsNames)
		if err != nil {
			return 0, nil, err
		}

		bestMovies = append(bestMovies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, bestMovies, nil
}

func (storage *MovieRepository) GetAllGenres() ([]string, error) {
	tx, err := storage.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()

	sqlStatement := `
		SELECT available_genres
		FROM mdb.meta
		ORDER BY version DESC
	`

	var genres []string
	err = tx.QueryRow(context.Background(), sqlStatement).Scan(&genres)
	if err == sql.ErrNoRows {
		return genres, nil
	}
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (storage *MovieRepository) GetMoviesByGenres(genres []string, startIndex int) (int, []*models.Movie, error) {
	tx, err := storage.db.Begin(context.Background())
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()

	var movies []*models.Movie

	sqlStatement := `
		SELECT COUNT(*)
		FROM mdb.movie
		WHERE genre && $1
	`
	var rowsCount int
	err = tx.QueryRow(context.Background(), sqlStatement, genres).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, movies, nil
	}
	if err != nil {
		return 0, nil, err
	}

	sqlStatement = `
        SELECT id, title, description, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie
		WHERE genre && $1
        ORDER BY rating DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := tx.Query(context.Background(), sqlStatement, genres, _const.MoviesPageSize, startIndex)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		movie := &models.Movie{}
		var id int
		var actorsNames []string
		err = rows.Scan(&id, &movie.Title, &movie.Description,
			&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
			&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
			&movie.Montage, &movie.Budget, &movie.Duration, &actorsNames, &movie.Poster,
			&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)
		if err != nil && err != sql.ErrNoRows {
			return 0, nil, err
		}

		movie.ID = strconv.Itoa(id)
		movie.Actors, err = storage.getActorsData(actorsNames)
		if err != nil {
			return 0, nil, err
		}

		movies = append(movies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, movies, nil
}

func (storage *MovieRepository) getActorsData(names []string) ([]models.ActorData, error) {
	tx, err := storage.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(context.Background())
		default:
			_ = tx.Rollback(context.Background())
		}
	}()

	var actors []models.ActorData

	sqlStatement := `
		SELECT id, name
		FROM mdb.actors
		WHERE name=$1
	`
	for _, name := range names {
		var actor models.ActorData
		err := tx.QueryRow(context.Background(), sqlStatement, name).Scan(&actor.ID, &actor.Name)
		if err != nil && err != sql.ErrNoRows {
			return []models.ActorData{}, err
		}

		actors = append(actors, actor)
	}

	return actors, nil
}
