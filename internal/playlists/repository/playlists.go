package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
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

type PlaylistsRepository struct {
	db PgxPoolIface
}

func NewPlaylistsRepository(database PgxPoolIface) *PlaylistsRepository {
	return &PlaylistsRepository{
		db: database,
	}
}

func (storage *PlaylistsRepository) CreatePlaylist(username string, playlistName string, isShared bool) error {
	sqlStatement := `
        INSERT INTO mdb.playlists (name, ownerName, isShared)
        VALUES ($1, $2, $3)
        RETURNING "id";
    `

	var newID int
	err := storage.db.
		QueryRow(context.Background(), sqlStatement,
			playlistName, username, isShared).
		Scan(&newID)

	if err != nil {
		return err
	}

	sqlStatement = `
        INSERT INTO mdb.playlistsWhoCanAdd (username, playlist_id)
        VALUES ($1, $2);
    `
	_, err = storage.db.
		Exec(context.Background(), sqlStatement,
			username, newID)

	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsRepository) GetPlaylistsInfo(username string, movieID int) ([]*models.PlaylistsInfo, error) {
	sqlStatement := `
        SELECT pl.id, pl.name, coalesce(plm.movie_id, -1) as movie_id
        FROM mdb.playlistsWhoCanAdd plwca
		LEFT JOIN mdb.playlistsMovies plm ON plwca.playlist_id = plm.playlist_id AND plm.movie_id = $1
		JOIN mdb.playlists pl ON plwca.playlist_id = pl.id
        WHERE plwca.username = $2
    `

	rows, err := storage.db.
		Query(context.Background(), sqlStatement, movieID, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlistsInfo []*models.PlaylistsInfo

	for rows.Next() {
		var newID int
		var playlistName string
		var newMovieID int

		err = rows.Scan(&newID, &playlistName, &newMovieID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		playlist := &models.PlaylistsInfo{}
		playlist.ID = strconv.Itoa(newID)
		playlist.Name = username

		if newMovieID == -1 {
			playlist.IsAdded = false
		} else {
			playlist.IsAdded = true
		}

		playlistsInfo = append(playlistsInfo, playlist)
	}

	return playlistsInfo, nil
}

func (storage *PlaylistsRepository) GetPlaylists(username string) ([]*models.Playlist, error) {
	sqlStatement := `
        SELECT pl.id, pl.name, json_object_agg(coalesce(m.id, -1), coalesce(m.title, '')) as kek
        FROM mdb.playlistsWhoCanAdd plwca
		LEFT JOIN mdb.playlistsMovies plm ON plwca.playlist_id = plm.playlist_id
		JOIN mdb.playlists pl ON plwca.playlist_id = pl.id
		LEFT JOIN mdb.movie m ON m.id = plm.movie_id
        WHERE plwca.username = $1
		GROUP BY pl.id
    `

	rows, err := storage.db.
		Query(context.Background(), sqlStatement, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlistsInfo []*models.Playlist

	for rows.Next() {
		var newID int
		var playlistName string
		movies := map[int]string{}

		err = rows.Scan(&newID, &playlistName, &movies)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		playlist := &models.Playlist{}
		playlist.ID = strconv.Itoa(newID)
		playlist.Name = playlistName

		for id, title := range movies {
			if id != -1 {
				newID := strconv.Itoa(id)
				playlist.Movies = append(playlist.Movies, models.MovieInPlaylist{newID, title, ""})
			}
		}

		playlistsInfo = append(playlistsInfo, playlist)
	}

	return playlistsInfo, nil
}

func (storage *PlaylistsRepository) CanUserUpdatePlaylist(username string, playlistID int) error {
	sqlStatement := `
	    SELECT username
		FROM mdb.playlists
	    WHERE id = $1 AND ownerName = $2;
	`

	_, err := storage.db.
		Query(context.Background(), sqlStatement, playlistID, username)

	if err != nil {
		return errors.New("user can't update playlist")
	}

	return nil
}

func (storage *PlaylistsRepository) DeleteAllUserFromPlaylist(username string, playlistID int) error {
	sqlStatement := `
	    DELETE
		FROM mdb.playlistsWhoCanAdd
	    WHERE playlist_id = $1 AND username <> $2;
	`

	_, err := storage.db.
		Query(context.Background(), sqlStatement, playlistID, username)

	if err != nil {
		return errors.New("can't delete user from playlist")
	}

	sqlStatement = `
	    DELETE
		FROM mdb.playlistsMovies
	    WHERE playlist_id = $1 AND addedBy <> $2;
	`

	_, err = storage.db.
		Query(context.Background(), sqlStatement, playlistID, username)

	if err != nil {
		return errors.New("can't delete user movie from playlist")
	}

	return nil
}

func (storage *PlaylistsRepository) UpdatePlaylist(username string, playlistID int, playlistName string, isShared bool) error {
	sqlStatement := `
	    UPDATE mdb.playlists
	    SET (name, isShared) = ($2, $3)
	    WHERE id=$1
		RETURNING isShared;
	`

	var NewIsShared bool

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, playlistID, playlistName, isShared).Scan(&NewIsShared)

	if err != nil {
		return errors.New("update playlist error")
	}

	if !NewIsShared {
		storage.DeleteAllUserFromPlaylist(username, playlistID)
	}

	return nil
}

func (storage *PlaylistsRepository) DeletePlaylist(playlistID int) error {
	sqlStatement := `
        DELETE FROM mdb.playlists
        WHERE id=$1;
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, playlistID)

	if err != nil {
		return errors.New("delete playlist error")
	}

	return nil
}

func (storage *PlaylistsRepository) CanUserUpdateMovieInPlaylist(username string, playlistID int) error {
	sqlStatement := `
	    SELECT addedBy
		FROM mdb.playlistsWhoCanAdd
	    WHERE playlist_id = $1 AND username = $2;
	`

	var newAddedBy string

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, playlistID, username).Scan(&newAddedBy)

	if err != nil {
		return errors.New("user can't add movie to playlist")
	}

	return nil
}

func (storage *PlaylistsRepository) AddMovieToPlaylist(username string, playlistID int, movieID int) error {
	sqlStatement := `
        INSERT INTO mdb.playlistsMovies (playlist_id, movie_id, addedBy)
        VALUES ($1, $2, $3)
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, playlistID, movieID, username)

	if err != nil {
		return errors.New("add movie to playlist error")
	}

	return nil
}

func (storage *PlaylistsRepository) DeleteMovieFromPlaylist(username string, playlistID int, movieID int) error {
	sqlStatement := `
        DELETE FROM mdb.playlistsMovies
        WHERE playlist_id = $1 AND movie_id = $2;
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, playlistID, movieID)

	if err != nil {
		return errors.New("delete movie from playlist error")
	}

	return nil
}

func (storage *PlaylistsRepository) CanUserUpdateUsersInPlaylist(username string, playlistID int) error {
	sqlStatement := `
	    SELECT ownerName
		FROM mdb.playlists
	    WHERE playlist_id = $1 AND ownerName = $2;
	`

	var newOwnerName string

	err := storage.db.
		QueryRow(context.Background(), sqlStatement, playlistID, username).Scan(&newOwnerName)

	if err != nil {
		return errors.New("user can't edit users in playlist")
	}

	return nil
}

func (storage *PlaylistsRepository) AddUserToPlaylist(username string, playlistID int, usernameToAdd string) error {
	sqlStatement := `
        INSERT INTO mdb.playlistsWhoCanAdd (username, playlist_id)
        VALUES ($1, $2)
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, usernameToAdd, playlistID)

	if err != nil {
		return errors.New("add user to playlist error")
	}

	return nil
}

func (storage *PlaylistsRepository) DeleteUserFromPlaylist(username string, playlistID int, usernameToDelete string) error {
	sqlStatement := `
        DELETE FROM mdb.playlistsWhoCanAdd
        WHERE username = $1 AND playlist_id = $2;
    `

	_, err := storage.db.
		Exec(context.Background(), sqlStatement, usernameToDelete, playlistID)

	if err != nil {
		return errors.New("delete user from playlist error")
	}

	return nil
}
