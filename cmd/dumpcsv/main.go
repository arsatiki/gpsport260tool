package main

import (
	"flag"
	"log"
	"os"
	"time"
	"fmt"
	"encoding/csv"

	"holux"
)

var (
	minTrackDuration = flag.Duration("drop", 5*60*time.Second,
		"Do not download messages shorter than specified duration")
)


func main() {
	flag.Parse()

	c, err := holux.Connect()

	if err != nil {
		log.Fatal(err)
	}
	c.Hello()
	defer c.Bye()

	index, err := c.GetIndex()
	if err != nil {
		log.Fatalf("Got error %v reading index, aborting", err)
	}
	
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{
		"index",
		"timestamp",
		"ridetime",
		"distance/m",
		"lat", "lon",
		"unk1",
		"flags",
		"fake_heading",
		"unk2",
	})
	known_flags := byte(0xcf)

	for k, track := range index {
		if !validTrack(track, *minTrackDuration) {
			continue
		}
		tzero := track.Time()

		points, err := c.GetTrack(track.Offset, track.Size)
		if err != nil {
			log.Print("Got error %v reading track %d", err, k)
			continue
		}

		for i := range points {
			fields := []string{
				fmt.Sprintf("%d", k),
				points[i].Time().Format(holux.TRKPTTIME),
				points[i].Time().Sub(tzero).String(),
				fmt.Sprintf("%d", points[i].Distance),
				fmt.Sprintf("%0.5f", points[i].Lat),
				fmt.Sprintf("%0.5f", points[i].Lon),
				fmt.Sprintf("%02x", points[i].Unk1),
				fmt.Sprintf("%02x", points[i].Flags & known_flags),
				fmt.Sprintf("%04x", points[i].Heading),
				fmt.Sprintf("%08x", points[i].Unk2),
			}
			w.Write(fields)
		}
		time.Sleep(300 * time.Millisecond)
	}
}


func validTrack(i holux.Index, minDuration time.Duration) bool {
	return i.Duration() > minDuration
}
