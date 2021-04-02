package localstorage

import (
	"context"
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
        SELECT id, title, description, voiceover, subtitles, quality, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               rating, rating_count
        FROM mdb.movie
        WHERE id=$1
    `

	idFilm, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = movieStorage.db.
		QueryRow(context.Background(), sqlStatement, idFilm).Scan(&idFilm,
		&movie.Title, &movie.Description, &movie.Voiceover, &movie.Subtitles, &movie.Quality,
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

func (movieStorage *MovieRepository) GetBestMovies(page, startIndex int) (int, []*models.Movie) {
	var bestMovies []*models.Movie

	sqlStatement := `
		SELECT COUNT(*) 
		FROM mdb.movie
	`

	var rowsCount int
	err := movieStorage.db.QueryRow(context.Background(), sqlStatement).Scan(&rowsCount)
	if err != nil {
		return 0, nil
	}

	sqlStatement = `
        SELECT id, title, description, voiceover, subtitles, quality, productionYear, country,
               genre, slogan, director, scriptwriter, producer, operator, composer,
               artist, montage, budget, duration, actors, poster, banner, trailerPreview,
               rating, rating_count
        FROM mdb.movie
        ORDER BY rating
        LIMIT $2 OFFSET $3
    `

	rows, err := movieStorage.db.Query(context.Background(), sqlStatement, _const.MoviesPageSize, startIndex)
	if err != nil {
		return 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		movie := &models.Movie{}
		var id int
		err = rows.Scan(&id, &movie.Title, &movie.Description, &movie.Voiceover, &movie.Subtitles, &movie.Quality,
			&movie.ProductionYear, &movie.Country, &movie.Genre, &movie.Slogan, &movie.Director,
			&movie.Scriptwriter, &movie.Producer, &movie.Operator, &movie.Composer, &movie.Artist,
			&movie.Montage, &movie.Budget, &movie.Duration, &movie.Actors, &movie.Poster,
			&movie.Banner, &movie.TrailerPreview, &movie.Rating, &movie.RatingCount)
		if err != nil {
			return 0, nil
		}
		movie.ID = strconv.Itoa(id)
		bestMovies = append(bestMovies, movie)
	}

	pagesNumber := int(math.Ceil(float64(rowsCount) / _const.MoviesPageSize))

	return pagesNumber, bestMovies
}
