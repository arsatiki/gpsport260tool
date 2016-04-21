package holux

import (
	"fmt"
	"math"
	"time"
)

const (
	y2000     = 946684800
	trackSize = 32
	indexSize = 64
	trkptTime = "2006-01-02 15:04:05Z07:00"
)

// FF 00 00 FF FF FF FF FF FF FF FF FF FF FF FF FF
//          01 jos favourite, muuten FF
// 52 08 B6 17 FD 0B 00 00 3F 04 00 00 00 00 00 00
// |--time---|             |distance-| |-offset--|
// 37 00 00 00 47 00 0C 00 2E 00 02 00 02 00 00 00
// |--size---| |smx| |sav| |cal|             HM HA
// 00 00 00 00 E6 00 00 00 02 00 00 00 00 00 00 00
// [pois] ?? ] [--ascent-] [-descent-]
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
	POIs     uint16 // Might be byte too.
	CADMax   byte
	CADAvg   byte
	Ascent   uint32 // meters
	Descent  uint32 // meters
	Unk2     [4]byte
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
	s := `Favorite: %v [Name: % 02x (%s)]
	Time: %v: Distance: %d m, Duration: %v
	Offset: %d points (%d B), Size: %d points (%d B)
	SPDMAX: %.1f km/h, SPDAVG: %.1f km/h, CAL: %d
	CO2 %.1f kg
	HRMMax: %d, HRMAvg: %d
	CADMax: %d, CADAvg: %d
	POIs: %d, Ascent: %d m, Descent: %d m
	Unknown: %s
	`
	return fmt.Sprintf(s, i.IsFavorite(), i.RawName, i.Name(),
		i.Time(), i.Distance, i.Duration(),
		i.Offset, i.Offset*32, i.Size, i.Size*32,
		float32(i.SpeedMax)/10, float32(i.SpeedAvg)/10, i.Calories,
		float32(i.CO2)/10,
		i.HRMMax, i.HRMAvg,
		i.CADMax, i.CADAvg,
		i.POIs, i.Ascent, i.Descent,
		i.UnknownFields())
}

type Trackpoint struct {
	RawTime            MTKTime
	Lat                float32 // North Positive
	Lon                float32 // East Positive
	GPSAltitude        int16
	Speed              uint16 // hm/h
	Unk1               byte   // Heading mod 255???
	Flags              TPFlag
	HR                 byte
	Cadence            byte
	BarometricAltitude int16
	Heading            uint16 // Not a heading tbh. Grade?
	Distance           uint32
	Unk2               uint32 // Never seen above 00
}

type Track []Trackpoint

func (t Track) NormalizeHR() {
	var adjustment int64

	for k, _ := range t {
		if k == 0 || t[k].HR == 0 || normalHRDelta(t[k-1], t[k]) {
			adjustment = 0
			continue
		}
		if adjustment == 0 {
			adjustment = int64(t[k-1].HR) - int64(t[k].HR)
		}
		t[k].HR = byte(int64(t[k].HR) + adjustment)
	}
}

func normalHRDelta(prev, curr Trackpoint) bool {
	Δt := curr.RawTime - prev.RawTime
	Δhr := math.Abs(float64(curr.HR) - float64(prev.HR))
	return Δt > 60 || Δhr < 30
}

func (t Trackpoint) IsPOI() bool {
	return t.Flags&flagTrkptPOI != 0
}

func (t Trackpoint) HasHR() bool {
	return t.Flags&flagTrkptHR != 0
}

func (t Trackpoint) HasCadence() bool {
	return t.Flags&flagTrkptCad != 0
}

func (t Trackpoint) Time() time.Time {
	return t.RawTime.Value()
}

func (t Trackpoint) UnknownFields() string {
	f := "Unk1 %02x|Flags %s|Unk2 %08x"
	return fmt.Sprintf(f, t.Unk1, t.Flags, t.Unk2)
}

// TODO: Add more fields, perhaps?
func (t Trackpoint) String() string {
	return fmt.Sprintf(`TRKPT: %s %s, %s
		Height (GPS): %d m, Speed: %.1f km/h,
		HR: %d, Cadence: %d, Alt (Baro): %d m, Heading: %d, Distance: %d m
		UnknownFields: %s
		`, t.Time().Format(trkptTime),
		fmtCoordinate(t.Lat, "N", "S"), fmtCoordinate(t.Lon, "E", "W"),
		t.GPSAltitude, float32(t.Speed)/10,
		t.HR, t.Cadence, t.BarometricAltitude, t.Heading, t.Distance,
		t.UnknownFields(),
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
	return time.Unix(int64(t)+y2000, 0)
}

// Trackpoint flag field bit masks
type TPFlag byte

const (
	flagTrkptUnk1 TPFlag = 1 << iota // 0x01 ? often set
	flagTrkptUnk2                    // 0x02 ? rarely set
	flagTrkptUnk3                    // 0x04 ? often set
	flagTrkptUnk4                    // 0x08 ? sometimes set, when 0x04 isn't
	flagTrkptPOI                     // 0x10 POI
	flagTrkptHR                      // 0x20 Heartrate present
	flagTrkptCad                     // 0x40 Cadence present
	flagTrkptUnk6                    // 0x80 ? never seen
)

func (f TPFlag) String() string {
	names := []string{
		"UK1",
		"UK2",
		"UK3",
		"UK4",
		"POI",
		"HRM",
		"CAD",
		"UK6",
	}

	seen := make([]string, 0, 8)

	for k, name := range names {
		if f&(1<<uint8(k)) != 0 {
			seen = append(seen, name)
		}
	}
	return fmt.Sprintf("% s", seen)
}
