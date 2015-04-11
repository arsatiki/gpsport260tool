package main

import (
	"database/sql"
	"holux"

	_ "github.com/mattn/go-sqlite3"
)

//go:generate awk -f generate.awk setup.sql
// TODO: Needs an insert test

var (
	TABLES = []string{
		"tracks", "trackpoints", "points_of_interest", "uploads",
	}
)

func initialize(db *sql.DB, err error) error {
	var name string

	if err != nil {
		return err
	}

	// TODO: supported by Snow Lep?
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	for _, table := range TABLES {
		err = db.QueryRow(NEEDS_INIT, table).Scan(&name)
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

func saveTrack(tx *sql.Tx, t holux.Index, err error) (int64, error) {
	if err != nil {
		return 0, err
	}

	trackname := sql.NullString{
		String: t.Name(),
		Valid:  t.IsNameSet(),
	}

	res, err := tx.Exec(INSERT["track"], t.Time(), trackname,
		t.Distance, t.RawDuration, t.UnknownFields())

	if err != nil {
		return 0, err
	}

	trackID, err := res.LastInsertId()
	return trackID, err
}

func savePoints(tx *sql.Tx, t holux.Track, trackID int64, err error) error {
	if err != nil {
		return err
	}

	insertPoint, err := tx.Prepare(INSERT["trackpoint"])

	if err != nil {
		return err
	}

	for _, point := range t {
		if !point.IsPOI() {
			hr := sql.NullInt64{
				Int64: int64(point.HR),
				Valid: point.HasHR(),
			}
			// TODO fix cadences
			cd := sql.NullInt64{
				Int64: 0,
				Valid: point.HasCadence(),
			}

			_, err = insertPoint.Exec(trackID, point.Time(),
				point.Lat, point.Lon, point.Alt, point.Height,
				hr, cd, point.UnknownFields())

		} else {
			_, err = tx.Exec(INSERT["POI"], trackID, point.Time(),
				point.Lat, point.Lon)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
