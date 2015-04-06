package main

import (
	"database/sql"
	"holux"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	c, err := holux.Connect()

	if err != nil {
		log.Fatal(err)
	}
	c.Hello()
	defer c.Bye()

	db, err := sql.Open("sqlite3", "tracks.db")
	err = initialize(db, err)
	if err != nil {
		log.Fatalf("error while initializing DB:", err)
	}
	defer db.Close()

	index, err := c.GetIndex()
	if err != nil {
		log.Fatalf("Got error %v, arborting", err)
	}

	for k, track := range index {
		points, err := c.GetTrack(track.Offset, track.Size)
		if err != nil {
			log.Fatal("Got error %v while reading track %d", err, k)
		}

		tx, err := db.Begin()
		trackID, err := saveTrack(tx, track, err)
		log.Println("TrackID is", trackID)
		err = savePoints(tx, points, trackID, err)

		if err != nil {
			tx.Rollback()
			log.Fatalf("Got error while writing track %d:", k, err)
		}
		tx.Commit()
	}
}
