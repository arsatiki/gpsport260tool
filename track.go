package holux

import (
	"time"
)

const (
	Y2000     = 946684800
	TRACKSIZE = 32
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
