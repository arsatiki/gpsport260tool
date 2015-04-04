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

// insert POIs

func saveTrack(tx *sql.Tx, t holux.Index) error {
	// TODO: Insert track index
	return nil
}

func savePoints(tx *sql.Tx, t holux.Track, index int64) error {
	writePoint, err := tx.Prepare(INSERT["trackpoint"])

	if err != nil {
		return err
	}

	for _, point := range t {
		hr := sql.NullInt64{Int64: int64(point.HR), Valid: point.HR != 0}
		cd := sql.NullInt64{Valid: false}

		_, err := writePoint.Exec(index, point.Time(),
			point.Lat, point.Lon, point.Height,
			hr, cd)

		if err != nil {
			return err
		}
	}

	return nil
}

func savePOIs(tx *sql.Tx, POIs []holux.Trackpoint) error {
	// TODO write pois
	return nil
}
