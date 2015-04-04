package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const TRACKSQL = `
	select time, latitude, longitude, elevation, heartrate, cadence
	from trackpoints
	where track = ?
	order by time
`

// TODO: Possible to hint at how many tracks are coming?
func GetTrackpoints(db *sql.DB, id int) ([]Trkpt, error) {
	var points []Trkpt

	rows, err := db.Query(TRACKSQL, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Trkpt
		var hr, cadence sql.NullInt64

		err = rows.Scan(&t.Time.Time,
			&t.Lat, &t.Lon, &t.Ele,
			&hr, &cadence)

		if err != nil {
			return points, err
		}
		t.HR = hr.Int64
		t.Cadence = cadence.Int64

		points = append(points, t)

	}
	return points, rows.Err()
}
