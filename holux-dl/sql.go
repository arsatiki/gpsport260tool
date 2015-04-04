package main

// Generated from setup.sql
const (
	NEEDS_INIT = `SELECT name FROM sqlite_master WHERE type='table' and name=?;`
)

var (
	CREATE_TABLE = map[string]string{
		"trackpoints":        `CREATE TABLE trackpoints ( time      TIMESTAMP NOT NULL, latitude  REAL NOT NULL, longitude REAL NOT NULL, elevation REAL NOT NULL, heartrate INTEGER, cadence   INTEGER, track     INTEGER NOT NULL REFERENCES tracks(ROWID) );`,
		"uploads":            `CREATE TABLE uploads ( track INTEGER NOT NULL REFERENCES tracks(ROWID), url   TEXT );`,
		"points_of_interest": `CREATE TABLE points_of_interest ( time        TIMESTAMP NOT NULL, latitude    REAL NOT NULL, longitude   REAL NOT NULL, description TEXT );`,
		"tracks":             `CREATE TABLE tracks ( time     TIMESTAMP NOT NULL, name     TEXT,  distance INTEGER NOT NULL,  duration INTEGER NOT NULL  );`,
	}
	INSERT = map[string]string{
		"POI":        `INSERT INTO points_of_interest(time, latitude, longitude) VALUES (?, ?, ?);`,
		"trackpoint": `INSERT INTO trackpoints(track, time, latitude, longitude, elevation, heartrate, cadence) VALUES (?, ?, ?);`,
		"track":      `INSERT INTO tracks(time, name, distance, duration)  VALUES (?, ?, ?, ?);`,
	}
)
