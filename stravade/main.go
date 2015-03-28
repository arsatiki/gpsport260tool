package main

import (
	_ "compress/gzip"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// TODO:
// - format time as Zulu time
// Support for more than one trkseg? Mebbe. Mebbe not.

type GPX struct {
	XMLName   xml.Name  `xml:"gpx"`
	XMLNS     string    `xml:"xmlns,attr"`
	XMLNSxsi  string    `xml:"xmlns:xsi,attr"`
	XMLSchema string    `xml:"xsi:schemaLocation,attr"`
	Creator   string    `xml:"creator,attr"`
	Version   string    `xml:"version,attr"`
	Time      time.Time `xml:"metadata>time"`

	Name   string  `xml:"trk>name"`
	Points []Trkpt `xml:"trk>trkseg>trkpt"`
}

type Trkpt struct {
	Lat     float32   `xml:"lat,attr"`
	Lon     float32   `xml:"lon,attr"`
	Ele     float32   `xml:"ele"`
	Time    time.Time `xml:"time"`
	HR      int64     `xml:"extensions>heartrate,omitempty"`
	Cadence int64     `xml:"extensions>cadence,omitempty"`
}

var point = Trkpt{
	Lat:     60.1732920,
	Lon:     24.9311040,
	Ele:     14.5,
	Time:    time.Now(),
	HR:      90,
	Cadence: 0,
}

func NewGPX(name string, t time.Time, pts []Trkpt) GPX {
	return GPX{
		XMLNS:     "http://www.topografix.com/GPX/1/1",
		XMLNSxsi:  "http://www.w3.org/2001/XMLSchema-instance",
		XMLSchema: "http://www.topografix.com/GPX/1/1",

		Creator: "Holux GPSSport 260 Pro with barometer",
		Version: "1.1",
		Time:    t,
		Name:    name,
		Points:  pts,
	}
}

func main() {
	doc := NewGPX("Joyride", time.Now(), []Trkpt{point})

	//dst := gzip.NewWriter(os.Stdout)
	dst := os.Stdout
	defer dst.Close()
	dst.Write([]byte(xml.Header))
	enc := xml.NewEncoder(dst)
	enc.Indent("", "    ")

	err := enc.Encode(doc)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
