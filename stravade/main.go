package main

import (
	"compress/gzip"
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	compressed = flag.Bool("z", false, "compress output with gzip")
	filename   = flag.String("f", "holux.db", "path to track database")
)

/*
var point = Trkpt{
	Lat:     60.1732920,
	Lon:     24.9311040,
	Ele:     14.5,
	Time:    GPXTime{time.Now()},
	HR:      90,
	Cadence: 0,
}*/

func main() {
	var dst io.Writer = os.Stdout

	flag.Parse()
	if *compressed {
		zdst := gzip.NewWriter(dst)
		defer zdst.Close()
		dst = io.Writer(zdst)
	}

	db, err := sql.Open("sqlite3", *filename)
	if err != nil {
		log.Fatal("err")
	}

	points, err := GetTrackpoints(db, 1)

	if err != nil {
		log.Fatal(err)
	}

	doc := NewGPX("SQLtest", time.Now(), points)

	dst.Write([]byte(xml.Header))
	enc := xml.NewEncoder(dst)
	enc.Indent("", "    ")

	err = enc.Encode(doc)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
