package localstorage

import (
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
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

func (movieStorage *MovieRepository) GetMovieByID(id string, username string) (*models.Movie, error) {
	var movie models.Movie

	sqlStatement := `
        SELECT mv.id, title, description, productionYear, country,
               array_agg(distinct (gs.name)) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(distinct (ac.name)) as actors, poster, banner, trailerPreview,
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

	isWatched := false
	if username != "" {
		sqlStatementWatched := `
			SELECT COUNT(*) as count
			FROM mdb.watched_movies
			WHERE user_login=$1 AND movie_id=$2
		`
		var rowsCount int
		err = movieStorage.db.QueryRow(context.Background(), sqlStatementWatched, username, id).Scan(&rowsCount)
		if err != nil {
			return nil, err
		}
		if rowsCount > 0 {
			isWatched = true
		}
	}

	movie.IsWatched = isWatched
	movie.ID = strconv.Itoa(idFilm)

	return &movie, nil
}

func (movieStorage *MovieRepository) GetBestMovies(startIndex int, username string) (int, []*models.Movie, error) {
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
               array_agg(distinct (gs.name)) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(distinct (ac.name)) as actors, poster, banner, trailerPreview,
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

		isWatched := false
		if username != "" {
			sqlStatementWatched := `
			SELECT COUNT(*) as count
			FROM mdb.watched_movies
			WHERE user_login=$1 AND movie_id=$2
		`
			var rowsCount int
			err = movieStorage.db.QueryRow(context.Background(), sqlStatementWatched, username, id).Scan(&rowsCount)
			if err != nil {
				return 0, nil, err
			}
			if rowsCount > 0 {
				isWatched = true
			}
		}

		movie.IsWatched = isWatched

		bestMovies = append(bestMovies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, bestMovies, nil
}

func (movieStorage *MovieRepository) GetAllGenres() ([]string, error) {
	sqlStatement := `
		SELECT name
		FROM mdb.genres
	`

	var genres []string
	rows, err := movieStorage.db.Query(context.Background(), sqlStatement)
	if err != nil {
		return genres, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre string
		err = rows.Scan(&genre)
		if err != nil {
			return []string{}, err
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

func (movieStorage *MovieRepository) GetMoviesByGenres(genres []string, startIndex int, username string) (int, []*models.Movie, error) {
	var movies []*models.Movie

	sqlStatement := `
		SELECT COUNT(*) as count
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
               array_agg(distinct (gs.name)) as genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, array_agg(distinct (ac.name)) as actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1) as rating, rating_count
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

		isWatched := false
		if username != "" {
			sqlStatementWatched := `
			SELECT COUNT(*) as count
			FROM mdb.watched_movies
			WHERE user_login=$1 AND movie_id=$2
		`
			var rowsCount int
			err = movieStorage.db.QueryRow(context.Background(), sqlStatementWatched, username, id).Scan(&rowsCount)
			if err != nil {
				return 0, nil, err
			}
			if rowsCount > 0 {
				isWatched = true
			}
		}

		movie.IsWatched = isWatched

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

func (movieStorage *MovieRepository) MarkWatched(username string, id int) error {
	sqlStatement := `
		INSERT INTO mdb.watched_movies VALUES($1, $2)
	`

	_, err := movieStorage.db.Exec(context.Background(), sqlStatement, username, id)
	if err != nil {
		return err
	}
	return nil
}

func (movieStorage *MovieRepository) MarkUnwatched(username string, id int) error {
	sqlStatement := `
		DELETE FROM mdb.watched_movies WHERE user_login=$1 AND movie_id=$2 
	`

	_, err := movieStorage.db.Exec(context.Background(), sqlStatement, username, id)
	if err != nil {
		return err
	}
	return nil
}

func (movieStorage *MovieRepository) SearchMovies(query string) ([]models.Movie, error) {
	sqlSearchMovies := `
		SELECT mv.id, title, poster, array_agg(distinct (gs.name)) as genre
		FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		WHERE lower(title) LIKE '%' || $1 || '%'
		GROUP BY mv.id
	`

	rows, err := movieStorage.db.Query(context.Background(), sqlSearchMovies, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		movie := models.Movie{}
		var id int

		err = rows.Scan(&id, &movie.Title, &movie.Poster, &movie.Genre)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		movie.ID = strconv.Itoa(id)
		movies = append(movies, movie)
	}

	return movies, nil
}

func (movieStorage *MovieRepository) GetSimilar(id string) ([]models.Movie, error) {
	// используем коэффициент Жаккара
	// обозначим X - количество пользователей, которые посмотрели фильм x;
	// Y - количество пользователей, которые посмотрели фильм y
	// sim(x, y) = |X && Y|/|X || Y|
	// Иначе: sim(x, y) = (кол-во пользователей, посмотревших оба фильма)/(кол-во пользователей, посмотревших или x, или y)

	sqlStatementSimilar := `
	SELECT mv.id, title, poster,
		CAST((SELECT COUNT(*) FROM (
			SELECT user_login
			FROM mdb.watched_movies
			WHERE movie_id = $1
			INTERSECT
			SELECT user_login
			FROM mdb.watched_movies
			WHERE movie_id = mv.id
		) AS watched_both_users) AS decimal)
		/
		(SELECT COUNT(*) FROM (
			SELECT user_login
			FROM mdb.watched_movies
			WHERE movie_id = $1
			UNION
			SELECT user_login
			FROM mdb.watched_movies
			WHERE movie_id = mv.id
		) AS watched_at_least_one_users)
		AS similarity_coefficient
	FROM mdb.movie mv
		WHERE mv.id != $1
		ORDER BY similarity_coefficient
		LIMIT $2
	`

	rows, err := movieStorage.db.Query(context.Background(), sqlStatementSimilar, id, _const.SimilarMoviesLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	var count int
	for rows.Next() {
		count++
		movie := models.Movie{}
		var id int
		var similarity float64

		err = rows.Scan(&id, &movie.Title, &movie.Poster, &similarity)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		movie.ID = strconv.Itoa(id)
		movies = append(movies, movie)
	}

	if count == _const.SimilarMoviesLimit {
		return movies, nil
	}
	// если данных о просмотрах мало, дополним выдачу фильмами того же жанра
	count = _const.SimilarMoviesLimit - count

	// id фильмов, которые уже рекомендовали, чтобы не повторять их
	idInt, _ := strconv.Atoi(id)
	alreadyAddedIDs := []int{idInt}
	for _, movie := range movies {
		idInt, _ = strconv.Atoi(movie.ID)
		alreadyAddedIDs = append(alreadyAddedIDs, idInt)
	}

	sqlStatementSameGenre := `
		SELECT mv.id, title, poster
		FROM mdb.movie mv
		JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
		JOIN mdb.genres gs ON mvgs.genre_id = gs.id
		WHERE mv.id != ALL($2)
		GROUP BY mv.id
		HAVING array_agg(array[gs.name])
		<@
		(
			SELECT array_agg(distinct (gs.name))
			FROM mdb.movie mv
			JOIN mdb.movie_genres mvgs ON mv.id = mvgs.movie_id
			JOIN mdb.genres gs ON mvgs.genre_id = gs.id
			WHERE mv.id = $1
		)
		LIMIT $3
	`

	rows, err = movieStorage.db.Query(context.Background(), sqlStatementSameGenre, id, alreadyAddedIDs, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		movie := models.Movie{}
		var id int

		err = rows.Scan(&id, &movie.Title, &movie.Poster)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		movie.ID = strconv.Itoa(id)
		movies = append(movies, movie)
	}

	return movies, err
}
