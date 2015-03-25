package holux

import (
	"bytes"
	"fmt"
	"time"
)

const (
	Y2000     = 946684800
	TRACKSIZE = 32
	TRKPTTIME = "2006-01-02 15:04:05Z07:00"
)

type Trackpoint struct {
	TimeMKT  uint32
	Lat      float32 // North Positive
	Lon      float32 // East Positive
	Height   uint16
	Speed    uint16
	_        byte
	Flags    byte
	HR       uint16
	Alt      uint16
	Heading  uint16
	Distance uint32
	_        uint32 // Cadence?
}

type Track []Trackpoint

func (t Trackpoint) IsPOI() bool {
	return t.Flags&0x10 == 1
}

func (t Trackpoint) Time() time.Time {
	return time.Unix(int64(t.TimeMKT)+Y2000, 0)
}

// TODO: Add more fields, perhaps?
func (t Trackpoint) String() string {
	var out bytes.Buffer

	fmt.Fprintf(&out, t.Time().Format(TRKPTTIME))
	fmt.Fprintf(&out, " %s, %s",
		fmtCoordinate(t.Lat, "N", "S"),
		fmtCoordinate(t.Lon, "E", "W"))
	return out.String()
}

func fmtCoordinate(v float32, pos, neg string) string {
	switch {
	case v > 0:
		return fmt.Sprintf("%0.5f °%s", v, pos)
	case v < 0:
		return fmt.Sprintf("%0.5f °%s", -v, neg)
	}
	return "0 °"
}
