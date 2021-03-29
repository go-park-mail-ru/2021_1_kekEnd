package localstorage

import (
    "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
    "github.com/jackc/pgx/v4/pgxpool"
    "strconv"
    "context"
)


type MovieRepository struct {
    db           *pgxpool.Pool
}

func NewMovieRepository(database *pgxpool.Pool) *MovieRepository {
    return &MovieRepository {
        db:           database,
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
               artist, montage, budget, duration, actors, poster, banner, trailerPreview
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
            &movie.Banner, &movie.TrailerPreview)

    if err != nil {
        return nil, err
    }

    movie.ID = strconv.Itoa(idFilm)

    return &movie, nil
}
