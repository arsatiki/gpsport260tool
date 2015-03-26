package holux

import (
	"bytes"
	"fmt"
	"time"
)

const (
	Y2000     = 946684800
	TRACKSIZE = 32
	INDEXSIZE = 64
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

// FF 00 00 FF FF FF FF FF FF FF FF FF FF FF FF FF
// 52 08 B6 17 FD 0B 00 00 3F 04 00 00 00 00 00 00
// |--time---|             |distance-| |-offset--|
// 37 00 00 00 47 00 0C 00 2E 00 02 00 02 00 00 00
// |--size---| |smx| |sav| |cal|             HM HA
// 00 00 00 00 E6 00 00 00 02 00 00 00 00 00 00 00
type Index struct {
	_        [4]byte  // TODO double check
	Name     [12]byte // offset=4, length=11????
	TimeMKT  uint32   // MKTTime
	Duration uint32   // seconds
	Distance uint32   // meters
	Offset   uint32   // LIST_MEM_START_OFFSET=28
	Size     uint32   // LIST_MEM_LENGTH_OFFSET=32 What's the diff with Length?
	SpeedMax uint16   // 35.6 km/h = 356.
	SpeedAvg uint16
	Calories uint16
	Unk1     [4]byte
	HRMMax   byte // BPM
	HRMAvg   byte
	Unk2     [16]byte
}

func (i Index) String() string {
	s := `Name: [% 02x] Time: %v: Distance: %d m, Duration: %v
	Offset: %d points (%d B), Size: %d points (%d B)
	SPDMAX: %.1f km/h, SPDAVG: %.1f km/h, CAL: %d
	[% 02x]
	HRMMax: %d, HRMAvg: %d
	[% 02x]
	`
	return fmt.Sprintf(s, i.Name,
		time.Unix(int64(i.TimeMKT)+Y2000, 0), i.Distance, time.Duration(i.Duration)*time.Second,
		i.Offset, i.Offset*32, i.Size, i.Size*32,
		float32(i.SpeedMax)/10, float32(i.SpeedAvg)/10, i.Calories,
		i.Unk1,
		i.HRMMax, i.HRMAvg,
		i.Unk2)
}
