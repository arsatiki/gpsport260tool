package main

// Generated from setup.sql
const (
	NEEDS_INIT = `SELECT name FROM sqlite_master WHERE type='table' and name=?;`
)

var (
	CREATE_TABLE = map[string]string{
		"trackpoints":        `CREATE TABLE trackpoints ( FOREIGN KEY(track) REFERENCES tracks(ROWID) NOT NULL ON DELETE CASCADE, time      TIMESTAMP NOT NULL, latitude  REAL NOT NULL, longitude REAL NOT NULL, elevation REAL NOT NULL, heartrate INTEGER, cadence   INTEGER );`,
		"uploads":            `CREATE TABLE uploads ( FOREIGN KEY(track) REFERENCES tracks(ROWID) NOT NULL ON DELETE CASCADE, url   TEXT );`,
		"points_of_interest": `CREATE TABLE points_of_interest ( FOREIGN KEY(track) REFERENCES tracks(ROWID) ON DELETE SET NULL, time        TIMESTAMP NOT NULL, latitude    REAL NOT NULL, longitude   REAL NOT NULL, description TEXT );`,
		"tracks":             `CREATE TABLE tracks ( time     TIMESTAMP NOT NULL, name     TEXT,  distance INTEGER NOT NULL,  duration INTEGER NOT NULL  );`,
	}
	INSERT = map[string]string{
		"POI":        `INSERT INTO points_of_interest(track, time, latitude, longitude) VALUES (?, ?, ?);`,
		"trackpoint": `INSERT INTO trackpoints(track, time, latitude, longitude, elevation, heartrate, cadence) VALUES (?, ?, ?);`,
		"track":      `INSERT INTO tracks(time, name, distance, duration)  VALUES (?, ?, ?, ?);`,
	}
)
