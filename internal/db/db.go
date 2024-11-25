package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

type Storage interface {
	DeleteSong(song Song) error
	AddSong(song SongData) error
	UpdateSong(song SongData) error
	Info(typeSort Sort) ([]SongData, error)
	GetTextSong(song Song) (string, error)
}

type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type SongData struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseData string `json:"releaseData"`
	TextSong    string `json:"textSong"`
	SongLink    string `json:"songLink"`
}

type Sort struct {
	TypeSort  string `json:"typeSort"`
	Direction string `json:"direction"`
}

type DB struct {
	db *sql.DB
}

func NewDataBase(dbname string) (Storage, error) {
	db, err := sql.Open("postgres", dbname)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (db *DB) DeleteSong(song Song) error {
	_, err := db.db.Exec("DELETE FROM music WHERE author=$1 AND song=$2", song.Group, song.Song)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddSong(song SongData) error {
	_, err := db.db.Exec("INSERT INTO music (author, song, releaseData, textSong, songLink) VALUES ($1, $2, $3, $4, $5)", song.Group, song.Song, song.ReleaseData, song.TextSong, song.SongLink)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateSong(song SongData) error {
	query := "UPDATE music SET"
	var params []interface{}
	var setClauses []string

	if song.ReleaseData != "" {
		setClauses = append(setClauses, " releaseData=$1")
		params = append(params, song.ReleaseData)
	}
	if song.TextSong != "" {
		setClauses = append(setClauses, " textSong=$2")
		params = append(params, song.TextSong)
	}
	if song.SongLink != "" {
		setClauses = append(setClauses, " songLink=$3")
		params = append(params, song.SongLink)
	}

	// Если нет полей для обновления, возвращаем ошибку
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Собираем финальный запрос
	query += strings.Join(setClauses, ", ")
	query += " WHERE author=$" + strconv.Itoa(len(params)+1) + " AND song=$" + strconv.Itoa(len(params)+2)

	params = append(params, song.Group, song.Song)

	_, err := db.db.Exec(query, params...)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Info(typeSort Sort) ([]SongData, error) {
	var songs []SongData
	rows, err := db.db.Query("SELECT * FROM music ORDER BY " + typeSort.TypeSort + " " + typeSort.Direction)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var song SongData
		if err := rows.Scan(&song.Group, &song.Song, &song.ReleaseData, &song.TextSong, &song.SongLink); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil

}

func (db *DB) GetTextSong(song Song) (string, error) {
	var text string
	row := db.db.QueryRow("SELECT textSong FROM music WHERE author=$1 AND song=$2", song.Group, song.Song)
	err := row.Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}
