package main

import (
	"encoding/xml"
	"time"

	"holux"
)

// TODO: Support for more than one trkseg? Mebbe. Mebbe not.
type GPX struct {
	XMLName   xml.Name `xml:"gpx"`
	XMLNS     string   `xml:"xmlns,attr"`
	XMLNSxsi  string   `xml:"xmlns:xsi,attr"`
	XMLSchema string   `xml:"xsi:schemaLocation,attr"`
	Creator   string   `xml:"creator,attr"`
	Version   string   `xml:"version,attr"`

	Time GPXTime `xml:"metadata>time"`
	Name string  `xml:"trk>name"`
	// TODO: Consider a nested struct here?
	Points []Trkpt `xml:"trk>trkseg>trkpt"`

	Repr string `xml:",comment"`
}

type Trkpt struct {
	Lat  float32 `xml:"lat,attr"`
	Lon  float32 `xml:"lon,attr"`
	Ele  float32 `xml:"ele"`
	Time GPXTime `xml:"time"`

	// Heartrate and cadence are stored in extensions
	// and may be empty.
	HR      byte `xml:"extensions>heartrate,omitempty"`
	Cadence byte `xml:"extensions>cadence,omitempty"`

	Repr string `xml:",comment"`
}

// TODO: Kanssa SQL Scan?
type GPXTime struct {
	time.Time
}

func (t GPXTime) MarshalText() ([]byte, error) {
	u := t.UTC()
	return []byte(u.Format(time.RFC3339)), nil
}

func NewGPX(name string, t time.Time, track holux.Track, repr string) GPX {
	pts := make([]Trkpt, len(track))
	for k, p := range track {
		pts[k] = Trkpt{
			Lat: p.Lat, Lon: p.Lon, Ele: float32(p.GPSAltitude),
			Time:    GPXTime{p.Time()},
			HR:      p.HR,
			Cadence: p.Cadence,
			Repr:    p.String(),
		}
	}

	return GPX{
		XMLNS:     "http://www.topografix.com/GPX/1/1",
		XMLNSxsi:  "http://www.w3.org/2001/XMLSchema-instance",
		XMLSchema: "http://www.topografix.com/GPX/1/1",

		Creator: "Holux GPSSport 260 Pro",
		Version: "1.1",
		Time:    GPXTime{t},
		Name:    name,
		Points:  pts,
		Repr:    repr,
	}
}
