package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"holux"
)

//go:generate awk -f generate.awk setup.sql

var (
	TABLES = []string{
		"tracks", "trackpoints", "points_of_interest", "uploads",
	}
)

func initialize(db *sql.DB) error {
	var name string

	for _, table := range TABLES {
		err := db.QueryRow(NEEDS_INIT, table).Scan(&name)
		switch {
		case err == nil:
			continue
		case err != sql.ErrNoRows:
			return err
		}

		_, err = db.Exec(CREATE_TABLE[table])

		if err != nil {
			return err
		}
	}
	return nil
}

func saveTrack(db *sql.DB, t holux.Track) error {
	// write metadata
	// insert tracks
	// insert POIs
	return nil
}
