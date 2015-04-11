-- NEEDS_INIT
SELECT name FROM sqlite_master WHERE type='table' and name=?;

-- CREATE_TABLE tracks
CREATE TABLE tracks (
    id   INTEGER PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    name TEXT, --XXX

    "distance/m" INTEGER NOT NULL, -- in meters
    "duration/s" INTEGER NOT NULL, -- in seconds

    unknown TEXT
);

-- CREATE_TABLE trackpoints
CREATE TABLE trackpoints (
    track     INTEGER NOT NULL,
    time      TIMESTAMP NOT NULL,
    latitude  REAL NOT NULL,
    longitude REAL NOT NULL,

    "elevation/m" REAL NOT NULL, -- corresponds to alt
    "height/m"    REAL NOT NULL, -- corresponds to height

    "heartrate/bpm" INTEGER,
    "cadence/rpm"   INTEGER,

    unknown TEXT,

    FOREIGN KEY(track) REFERENCES tracks(id)
                       ON DELETE CASCADE
);
-- CREATE_INDEX trackpoints
CREATE INDEX trackpointindex ON trackpoints(track);

-- CREATE_TABLE uploads
CREATE TABLE uploads (
    track INTEGER NOT NULL,
    url   TEXT,
    FOREIGN KEY(track) REFERENCES tracks(id)
                       ON DELETE CASCADE
);
-- CREATE_INDEX uploads
CREATE INDEX uploadsindex ON uploads(track);

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
-- CREATE_INDEX points_of_interest
CREATE INDEX POIsindex ON points_of_interest(track);

-- INSERT track
INSERT INTO tracks(time, name, "distance/m", "duration/s", unknown)
       VALUES (?, ?, ?, ?, ?);

-- INSERT POI
INSERT INTO points_of_interest(track, time, latitude, longitude)
       VALUES (?, ?, ?, ?);

-- INSERT trackpoint
INSERT INTO trackpoints
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

