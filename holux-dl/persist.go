package main

import (
	"holux"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate awk -f generate.awk setup.sql

func initialize(db *sql.DB) error {
	// check if exists.
	// prepare all insert statements?
	return nil
}

func saveTrack(db *sql.DB, t holux.Track) error {
	// write metadata
	// insert tracks
	// insert POIs
	return nil
}

