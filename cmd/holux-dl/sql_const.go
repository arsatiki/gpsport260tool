package main

// Generated from setup.sql
const (
	NEEDS_INIT = `SELECT name FROM sqlite_master WHERE name=?;`
)

var (
	CREATE_INDEX = map[string]string{
		"uploads":            `CREATE INDEX uploadsindex ON uploads(track);`,
		"points_of_interest": `CREATE INDEX POIsindex ON points_of_interest(track);`,
		"trackpoints":        `CREATE INDEX trackpointindex ON trackpoints(track);`,
	}
	CREATE_TABLE = map[string]string{
		"trackpoints":        `CREATE TABLE trackpoints ( track     INTEGER NOT NULL, time      TIMESTAMP NOT NULL, latitude  REAL NOT NULL, longitude REAL NOT NULL, "elevation/m" REAL NOT NULL,  "height/m"    REAL NOT NULL,  "heartrate/bpm" INTEGER, "cadence/rpm"   INTEGER, FOREIGN KEY(track) REFERENCES tracks(id) ON DELETE CASCADE );`,
		"uploads":            `CREATE TABLE uploads ( track INTEGER NOT NULL, url   TEXT, FOREIGN KEY(track) REFERENCES tracks(id) ON DELETE CASCADE );`,
		"points_of_interest": `CREATE TABLE points_of_interest ( track       INTEGER, time        TIMESTAMP NOT NULL, latitude    REAL NOT NULL, longitude   REAL NOT NULL, description TEXT, FOREIGN KEY(track) REFERENCES tracks(id) ON DELETE SET NULL );`,
		"tracks":             `CREATE TABLE tracks ( id   INTEGER PRIMARY KEY, time TIMESTAMP NOT NULL, name TEXT,  "distance/m" INTEGER NOT NULL,  "duration/s" INTEGER NOT NULL,  unknown TEXT );`,
	}
	INSERT = map[string]string{
		"POI":        `INSERT INTO points_of_interest(track, time, latitude, longitude) VALUES (?, ?, ?, ?);`,
		"trackpoint": `INSERT INTO trackpoints(track, time, latitude, longitude, "elevation/m", "height/m", "heartrate/bpm", "cadence/rpm") VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		"track":      `INSERT INTO tracks(time, name, "distance/m", "duration/s", unknown) VALUES (?, ?, ?, ?, ?);`,
	}
)
