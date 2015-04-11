package holux

import (
	"fmt"
	"time"
)

const (
	Y2000     = 946684800
	TRACKSIZE = 32
	INDEXSIZE = 64
	TRKPTTIME = "2006-01-02 15:04:05Z07:00"
)

// FF 00 00 FF FF FF FF FF FF FF FF FF FF FF FF FF
//          01 jos favourite, muuten FF
// 52 08 B6 17 FD 0B 00 00 3F 04 00 00 00 00 00 00
// |--time---|             |distance-| |-offset--|
// 37 00 00 00 47 00 0C 00 2E 00 02 00 02 00 00 00
// |--size---| |smx| |sav| |cal|             HM HA
// 00 00 00 00 E6 00 00 00 02 00 00 00 00 00 00 00
// [pois] ?? ]
type Index struct {
	F00         [3]byte // TODO double check
	UnkFlag     byte    // FF when not favourite, 01 when fav
	RawName     [10]byte
	Unk         [2]byte // First byte can be \0 for C strings
	RawTime     MTKTime // MKTTime
	RawDuration uint32  // seconds
	Distance    uint32  // meters

	Offset uint32 // LIST_MEM_START_OFFSET=28
	Size   uint32 // LIST_MEM_LENGTH_OFFSET=32

	SpeedMax uint16 // 35.6 km/h = 356.
	SpeedAvg uint16
	Calories uint16
	Unk1     [2]byte
	CO2      uint16 // hectograms. 1 hg = 100 g
	HRMMax   byte   // BPM
	HRMAvg   byte
	POIs     byte // Can be uint16 or 32 as well.
	Unk2     [15]byte
}

func (i Index) Name() string {
	if i.IsNameSet() {
		return string(i.RawName[:])
	}
	return ""
}

func (i Index) IsNameSet() bool {
	return i.RawName[0] != 0xff
}

func (i Index) Duration() time.Duration {
	return time.Duration(i.RawDuration) * time.Second
}

func (i Index) Time() time.Time {
	return i.RawTime.Value()
}

func (i Index) IsFavorite() bool {
	return i.UnkFlag == 0x01
}

func (i Index) UnknownFields() string {
	f := "F00 %02x|UnkFlag %02x|Unk %02x|Unk1 %02x|Unk2 %02x"
	return fmt.Sprintf(f, i.F00, i.UnkFlag, i.Unk, i.Unk1, i.Unk2)
}

func (i Index) String() string {
	s := `[FF0000: %02x] Favorite: %v [Name: % 02x (%s)] [Unk: %02x]
	Time: %v: Distance: %d m, Duration: %v
	Offset: %d points (%d B), Size: %d points (%d B)
	SPDMAX: %.1f km/h, SPDAVG: %.1f km/h, CAL: %d
	[% 02x]
	CO2 %.1f kg
	HRMMax: %d, HRMAvg: %d
	[% 02x] (starts with # of POIs, length?)
	`
	return fmt.Sprintf(s, i.F00, i.IsFavorite(), i.RawName, i.Name(), i.Unk,
		i.Time(), i.Distance, i.Duration(),
		i.Offset, i.Offset*32, i.Size, i.Size*32,
		float32(i.SpeedMax)/10, float32(i.SpeedAvg)/10, i.Calories,
		i.Unk1,
		float32(i.CO2)/10,
		i.HRMMax, i.HRMAvg,
		i.Unk2)
}

// Trackpoint flag field bit masks
const (
	FLAG_TRKPT_UNK1 = 1 << iota // 0x01 ? often set
	FLAG_TRKPT_UNK2             // 0x02 ? rarely set
	FLAG_TRKPT_UNK3             // 0x04 ? often set
	FLAG_TRKPT_UNK4             // 0x08 ? sometimes set, when 0x04 isn't
	FLAG_TRKPT_POI              // 0x10 POI
	FLAG_TRKPT_HR               // 0x20 Heartrate present
	FLAG_TRKPT_UNK5             // 0x40 ? never seen
	FLAG_TRKPT_UNK6             // 0x80 ? never seen
	// Last two are probably cadence and speed
)

type Trackpoint struct {
	RawTime  MTKTime
	Lat      float32 // North Positive
	Lon      float32 // East Positive
	Height   int16
	Speed    uint16
	Unk1     byte
	Flags    byte
	HR       uint16
	Alt      int16
	Heading  uint16
	Distance uint32
	Unk2     uint32 // Cadence?
}

type Track []Trackpoint

func (t Trackpoint) IsPOI() bool {
	return t.Flags&FLAG_TRKPT_POI != 0
}

func (t Trackpoint) HasHR() bool {
	return t.Flags&FLAG_TRKPT_HR != 0
}

// TODO: Find out which flag controls cadence
func (t Trackpoint) HasCadence() bool {
	return false
}

func (t Trackpoint) Time() time.Time {
	return t.RawTime.Value()
}

func (t Trackpoint) UnknownFields() string {
	f := "Unk1 %02x|Flags %02x|Unk2 %02x"
	return fmt.Sprintf(f, t.Unk1, t.Flags, t.Unk2)
}

// TODO: Add more fields, perhaps?
func (t Trackpoint) String() string {
	return fmt.Sprintf(`TRKPT: %s %s, %s
		Height: %d m, Speed: %.1f m/s, Flags: %02x,
		HR: %d, Alt: %d m, Heading: %d, Distance: %d m
		
		`, t.Time().Format(TRKPTTIME),
		fmtCoordinate(t.Lat, "N", "S"), fmtCoordinate(t.Lon, "E", "W"),
		t.Height, float32(t.Speed)/10, t.Flags,
		t.HR, t.Alt, t.Heading, t.Distance,
	)
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

type MTKTime uint32

func (t MTKTime) Value() time.Time {
	return time.Unix(int64(t)+Y2000, 0)
}
