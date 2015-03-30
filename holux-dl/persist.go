package main

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

const (
	DB_INIT = []string{
		"create table xxx",
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
