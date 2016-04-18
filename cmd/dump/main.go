package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"time"

	"holux"
)

var (
	periodLength = flag.Duration("period", 24*time.Hour,
		"Download items from the last period")
	minTrackDuration = flag.Duration("drop", 60*time.Second,
		"Do not download messages shorter than specified duration")
	//clearDeviceLog = flag.Bool("clear", false,
	//	"Clear log on device after succesful transfer")
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
		log.Fatalf("Got error %v, aborting", err)
	}

	for k, track := range index {
		if !validTrack(track, *minTrackDuration, *periodLength) {
			continue
		}

		points, err := c.GetTrack(track.Offset, track.Size)
		if err != nil {
			log.Fatal("Got error %v while reading track %d", err, k)
		}

		points.NormalizeHR()

		dst, err := os.Create(nameForTrack(track))
		if err != nil {
			log.Fatalf("create error: %v\n", err)
		}

		doc := NewGPX(nameForUpload(track), track.Time(), points, track.String())

		dst.Write([]byte(xml.Header))
		enc := xml.NewEncoder(dst)
		enc.Indent("", "    ")

		err = enc.Encode(doc)

		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

	}

	//if *clearDeviceLog {
	//	c.ClearLog()
	//}
}

func nameForTrack(i holux.Index) string {
	return i.Time().Format("ride 2006-01-02 150405.gpx")
}

func nameForUpload(i holux.Index) string {
	if i.IsNameSet() {
		return i.Name()
	}

	return i.Time().Format("2006-01-02 15:04:05 -0700")
}

func validTrack(i holux.Index, minDuration time.Duration, since time.Duration) bool {
	return i.Duration() > minDuration && time.Since(i.Time()) < since
}

