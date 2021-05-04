package localstorage

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
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

type MovieRepository struct {
	db PgxPoolIface
}

func NewMovieRepository(database PgxPoolIface) *MovieRepository {
	return &MovieRepository{
		db: database,
	}
}

func (movieStorage *MovieRepository) CreateMovie(movie *models.Movie) error {
	return nil
}

func (movieStorage *MovieRepository) GetMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie

	sqlStatement := `
        SELECT mv.id, title, description, productionYear, country,
               array_agg(array[gs.name]) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(DISTINCT array[ac.name]) as actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1) AS rating, rating_count
        FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		JOIN mdb.movie_actors mvac ON mv.id = mvac.movie_id
		JOIN mdb.actors ac ON mvac.actor_id = ac.id
		WHERE mv.id = $1
		GROUP BY mv.id
    `

	idFilm, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var actorsNames []string
	err = movieStorage.db.
		QueryRow(context.Background(), sqlStatement, idFilm).Scan(&idFilm,
		&movie.Title, &movie.Description,
		&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
		&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
		&movie.Montage, &movie.Budget, &movie.Duration, &actorsNames, &movie.Poster,
		&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)

	if err != nil {
		return nil, err
	}

	movie.Actors, err = movieStorage.getActorsData(actorsNames)
	if err != nil {
		return nil, err
	}

	movie.ID = strconv.Itoa(idFilm)

	return &movie, nil
}

func (movieStorage *MovieRepository) GetBestMovies(startIndex int) (int, []*models.Movie, error) {
	var bestMovies []*models.Movie

	sqlStatement := `
		SELECT movies_count
		FROM mdb.meta
		ORDER BY version DESC
	`

	var rowsCount int
	err := movieStorage.db.QueryRow(context.Background(), sqlStatement).Scan(&rowsCount)
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
        SELECT mv.id, title, description, productionYear, country,
               array_agg(array[gs.name]) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(DISTINCT array[ac.name]) as actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		JOIN mdb.movie_actors mvac ON mv.id = mvac.movie_id
		JOIN mdb.actors ac ON mvac.actor_id = ac.id
        GROUP BY mv.id
		ORDER BY rating DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := movieStorage.db.Query(context.Background(), sqlStatement, _const.MoviesPageSize, startIndex)
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
		movie.Actors, err = movieStorage.getActorsData(actorsNames)
		if err != nil {
			return 0, nil, err
		}

		bestMovies = append(bestMovies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, bestMovies, nil
}

func (movieStorage *MovieRepository) GetAllGenres() ([]string, error) {
	sqlStatement := `
		SELECT available_genres
		FROM mdb.meta
		ORDER BY version DESC
	`

	var genres []string
	err := movieStorage.db.QueryRow(context.Background(), sqlStatement).Scan(&genres)
	if err == sql.ErrNoRows {
		return genres, nil
	}
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (movieStorage *MovieRepository) GetMoviesByGenres(genres []string, startIndex int) (int, []*models.Movie, error) {
	var movies []*models.Movie

	sqlStatement := `
		SELECT COUNT(*)
        FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		GROUP BY mv.id
		HAVING array_agg(array[gs.name]) && $1
	`
	var rowsCount int
	err := movieStorage.db.QueryRow(context.Background(), sqlStatement, genres).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, movies, nil
	}
	if err != nil {
		return 0, nil, err
	}

	sqlStatement = `
        SELECT mv.id, title, description, productionYear, country,
               array_agg(array[gs.name]) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(DISTINCT array[ac.name]) as actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		JOIN mdb.movie_actors mvac ON mv.id = mvac.movie_id
		JOIN mdb.actors ac ON mvac.actor_id = ac.id
		GROUP BY mv.id
		HAVING array_agg(array[gs.name]) && $1
        ORDER BY rating DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := movieStorage.db.Query(context.Background(), sqlStatement, genres, _const.MoviesPageSize, startIndex)
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
		movie.Actors, err = movieStorage.getActorsData(actorsNames)
		if err != nil {
			return 0, nil, err
		}

		movies = append(movies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, movies, nil
}

func (movieStorage *MovieRepository) getActorsData(names []string) ([]models.ActorData, error) {
	var actors []models.ActorData

	sqlStatement := `
		SELECT id, name
		FROM mdb.actors
		WHERE name=$1
	`
	for _, name := range names {
		var actor models.ActorData
		err := movieStorage.db.QueryRow(context.Background(), sqlStatement, name).Scan(&actor.ID, &actor.Name)
		if err != nil && err != sql.ErrNoRows {
			return []models.ActorData{}, err
		}

		actors = append(actors, actor)
	}

	return actors, nil
}
