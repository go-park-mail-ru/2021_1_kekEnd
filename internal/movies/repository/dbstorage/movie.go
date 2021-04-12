package localstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"strconv"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(database *pgxpool.Pool) *MovieRepository {
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

	err = movieStorage.db.
		QueryRow(context.Background(), sqlStatement, idFilm).Scan(&idFilm,
		&movie.Title, &movie.Description,
		&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
		&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
		&movie.Montage, &movie.Budget, &movie.Duration, &movie.Actors, &movie.Poster,
		&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)

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
        SELECT id, title, description, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie
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
		err = rows.Scan(&id, &movie.Title, &movie.Description,
			&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
			&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
			&movie.Montage, &movie.Budget, &movie.Duration, &movie.Actors, &movie.Poster,
			&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)
		if err != nil && err != sql.ErrNoRows {
			return 0, nil, err
		}
		movie.ID = strconv.Itoa(id)
		bestMovies = append(bestMovies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, bestMovies, nil
}

func (movieStorage *MovieRepository) GetMoviesByGenres(genres []string, startIndex int) (int, []*models.Movie, error) {
	var movies []*models.Movie

	sqlStatement := `
		SELECT COUNT(*)
		FROM mdb.movie
		WHERE genre @> $1
	`
	var rowsCount int
	err := movieStorage.db.QueryRow(context.Background(), sqlStatement, genres).Scan(&rowsCount)
	if err == sql.ErrNoRows {
		return 0, movies, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}

	sqlStatement = `
        SELECT id, title, description, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               ROUND(CAST(rating AS numeric), 1), rating_count
        FROM mdb.movie
		WHERE genre @> $1
        ORDER BY rating DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := movieStorage.db.Query(context.Background(), sqlStatement, genres, _const.MoviesPageSize, startIndex)
	if err != nil {
		fmt.Println(err)
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		movie := &models.Movie{}
		var id int
		err = rows.Scan(&id, &movie.Title, &movie.Description,
			&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
			&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
			&movie.Montage, &movie.Budget, &movie.Duration, &movie.Actors, &movie.Poster,
			&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)
		if err != nil && err != sql.ErrNoRows {
			return 0, nil, err
		}
		movie.ID = strconv.Itoa(id)
		movies = append(movies, movie)
	}

	fmt.Println(rowsCount)
	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, movies, nil
}
