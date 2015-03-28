package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

const EXAMPLE = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx 
	creator="strava.com iPhone" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
 <metadata>
  <time>2015-03-27T14:30:01Z</time>
 </metadata>
 <trk>
  <name>Afternoon Ride</name>
  <trkseg>
   <trkpt lat="60.1732920" lon="24.9311040">
    <ele>14.5</ele>
    <time>2015-03-27T14:40:22Z</time>
    <extensions>
     <gpxtpx:TrackPointExtension>
      <gpxtpx:hr>141</gpxtpx:hr>
     </gpxtpx:TrackPointExtension>
    </extensions>
   </trkpt>
`

var point = Trkpt{
	Lat:     60.1732920,
	Lon:     24.9311040,
	Ele:     14.5,
	Time:    time.Now(),
	HR:      90,
	Cadence: 0,
}
var doc = GPX{
	Creator: "Hocus pocus",
	Version: "1.1",
	Time:    time.Now(),
	Track: Trk{
		Name:     "Joyride",
		Segments: []Trkseg{{Points: []Trkpt{point}}},
	},
}

type GPX struct {
	XMLName xml.Name  `xml:"gpx"`
	Creator string    `xml:"creator,attr"`
	Version string    `xml:"version,attr"`
	Time    time.Time `xml:"metadata>time"`
	Track   Trk       `xml:"trk"`
}

type Trk struct {
	Name     string   `xml:"name"`
	Segments []Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Points []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Lat     float32   `xml:"lat,attr"`
	Lon     float32   `xml:"lon,attr"`
	Ele     float32   `xml:"ele"`
	Time    time.Time `xml:"time"`
	HR      int64     `xml:"extensions>heartrate,omitempty"`
	Cadence int64     `xml:"extensions>cadence,omitempty"`
}

func main() {
	dst := gzip.NewWriter(os.Stdout)
	defer dst.Close()
	dst.Write([]byte(xml.Header))
	enc := xml.NewEncoder(dst)
	enc.Indent("", "    ")

	err := enc.Encode(doc)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
