package main

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

const (
	DB_INIT = []string{
	`CREATE TABLE trackpoints (
		time      TIMESTAMP NOT NULL,
		latitude  REAL NOT NULL,
		longitude REAL NOT NULL,
		elevation REAL NOT NULL,
		heartrate INTEGER,
		cadence   INTEGER,
		track     INTEGER NOT NULL REFERENCES tracks(ROWID)
	)`,
		"create table yyy",
		"create table zzz",
	}
	WRITE_TRACK = "insert into ... values ..."
	WRITE_POI   = "insert into ... values ..."
	WRITE_POINT = "insert into ... values ..."
)

func initialize(db *sql.DB) error {
	return nil
}

func saveTrack(db *sql.DB, t Track) error {
	// write metadata
	// insert tracks
	// insert POIs
}

// TODO: Prepare queries
