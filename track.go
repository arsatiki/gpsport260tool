package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

const Y2000 = 946684800

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
	_	 uint32 // Cadence?
}

func (t Trackpoint) IsPOI() bool {
	return t.Flags & 0x10 == 1
}

func (t Trackpoint) Time() time.Time {
	return time.Unix(t.TimeMKT + Y2000, 0)
}

func main() {
	var t Trackpoint
	data := []byte{
		0xc2, 0x0d, 0xb0, 0x1b, 0x72, 0x0a, 0x60, 0x42,
		0xbb, 0x49, 0x17, 0x42, 0xb3, 0x00, 0x12, 0x00,
		0xea, 0x05, 0x00, 0x00, 0xa9, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	fmt.Printf("%v", t)
}
