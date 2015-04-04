-- NEEDS_INIT
SELECT name FROM sqlite_master WHERE type='table' and name=?;

-- CREATE_TABLE tracks
CREATE TABLE tracks (
    time     TIMESTAMP NOT NULL,
    name     TEXT, --XXX
    distance INTEGER NOT NULL, -- in meters
    duration INTEGER NOT NULL -- in seconds,
);

-- CREATE_TABLE trackpoints
CREATE TABLE trackpoints (
    time      TIMESTAMP NOT NULL,
    latitude  REAL NOT NULL,
    longitude REAL NOT NULL,
    elevation REAL NOT NULL,
    heartrate INTEGER,
    cadence   INTEGER,
    track     INTEGER NOT NULL REFERENCES tracks(ROWID)
);

-- CREATE_TABLE uploads
CREATE TABLE uploads (
    track INTEGER NOT NULL REFERENCES tracks(ROWID),
    url   TEXT
);

-- CREATE_TABLE points_of_interest
CREATE TABLE points_of_interest (
    time        TIMESTAMP NOT NULL,
    latitude    REAL NOT NULL,
    longitude   REAL NOT NULL,
    description TEXT
);

-- INSERT track
INSERT INTO tracks(time, name, distance, duration) 
       VALUES (?, ?, ?, ?);

-- INSERT POI
INSERT INTO points_of_interest(time, latitude, longitude)
       VALUES (?, ?, ?);

-- INSERT trackpoint
INSERT INTO trackpoints(track, time, latitude, longitude, elevation,
                        heartrate, cadence)
       VALUES (?, ?, ?);

