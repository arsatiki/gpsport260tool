package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const TRACKSQL = `
	select time, latitude, longitude, elevation, heartrate, cadence
	from points
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
		rows.Scan(&t.Time.Time,
			&t.Lat, &t.Lon, &t.Ele,
			&t.HR, &t.Cadence)
		points = append(points, t)
	}
	return points, rows.Err()
}
