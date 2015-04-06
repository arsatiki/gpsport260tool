-- NEEDS_INIT
SELECT name FROM sqlite_master WHERE type='table' and name=?;

-- CREATE_TABLE tracks
CREATE TABLE tracks (
    id       INTEGER PRIMARY KEY,
    time     TIMESTAMP NOT NULL,
    name     TEXT, --XXX
    distance INTEGER NOT NULL, -- in meters
    duration INTEGER NOT NULL -- in seconds,
);

-- CREATE_TABLE trackpoints
CREATE TABLE trackpoints (
    track     INTEGER NOT NULL,
    time      TIMESTAMP NOT NULL,
    latitude  REAL NOT NULL,
    longitude REAL NOT NULL,
    elevation REAL NOT NULL,
    heartrate INTEGER,
    cadence   INTEGER,
    FOREIGN KEY(track) REFERENCES tracks(id)
                       ON DELETE CASCADE
);

-- CREATE_TABLE uploads
CREATE TABLE uploads (
    track INTEGER NOT NULL,
    url   TEXT,
    FOREIGN KEY(track) REFERENCES tracks(id)
                       ON DELETE CASCADE
);

-- CREATE_TABLE points_of_interest
CREATE TABLE points_of_interest (
    track       INTEGER,
    time        TIMESTAMP NOT NULL,
    latitude    REAL NOT NULL,
    longitude   REAL NOT NULL,
    description TEXT,
    FOREIGN KEY(track) REFERENCES tracks(id)
                       ON DELETE SET NULL
);

-- INSERT track
INSERT INTO tracks(time, name, distance, duration) 
       VALUES (?, ?, ?, ?);

-- INSERT POI
INSERT INTO points_of_interest(track, time, latitude, longitude)
       VALUES (?, ?, ?);

-- INSERT trackpoint
INSERT INTO trackpoints(track, time, latitude, longitude, elevation,
                        heartrate, cadence)
       VALUES (?, ?, ?);

