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
// |--size---|
// 00 00 00 00 E6 00 00 00 02 00 00 00 00 00 00 00
type Index struct {
	_        [4]byte  // TODO double check
	Name     [12]byte // offset=4, length=11????
	TimeMKT  uint32   // MKTTime
	Duration uint32   // seconds
	Distance uint32   // meters
	Offset   uint32   // LIST_MEM_START_OFFSET=28
	Size     uint32   // LIST_MEM_LENGTH_OFFSET=32 What's the diff with Length?
	Unk	 [12]byte
	Int0     uint32
	Int1     uint32
	Int2     uint32
	Int3     uint32
}

func (i Index) String() string {
	return fmt.Sprintf("Time: %v: Distance: %d, Duration: %v, Offset: %d tracks (%d B), Size: %d tracks (%d B)\nInts: %d, %d, %d, %d\n% 02x",
		time.Unix(int64(i.TimeMKT)+Y2000, 0), i.Distance, time.Duration(i.Duration)*time.Second,
		i.Offset, i.Offset*32, i.Size, i.Size*32,
		i.Int0, i.Int1, i.Int2, i.Int3,
		i.Unk)
}

/*
Time: 2012-08-09 05:34:10 +0000 UTC: Distance: 1087, Offset: 0 tracks/0 bytes, Size: 55 tracks/1760 bytes
Time: 2015-03-07 17:06:07 +0000 UTC: Distance: 175, Offset: 56 tracks/1792 bytes, Size: 11 tracks/352 bytes
Time: 2015-03-07 21:46:03 +0000 UTC: Distance: 1020, Offset: 67 tracks/2144 bytes, Size: 46 tracks/1472 bytes
Time: 2015-03-08 10:53:39 +0000 UTC: Distance: 1379, Offset: 113 tracks/3616 bytes, Size: 99 tracks/3168 bytes
Time: 2015-03-08 11:55:06 +0000 UTC: Distance: 926, Offset: 212 tracks/6784 bytes, Size: 19 tracks/608 bytes
Time: 2015-03-09 05:27:40 +0000 UTC: Distance: 12505, Offset: 231 tracks/7392 bytes, Size: 459 tracks/14688 bytes
*/
