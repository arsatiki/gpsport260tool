package holux

import (
	"testing"
)

var (
	data = []byte{
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x52, 0x08, 0xB6, 0x17, 0xFD, 0x0B, 0x00, 0x00, 0x3F, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x37, 0x00, 0x00, 0x00, 0x47, 0x00, 0x0C, 0x00, 0x2E, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xE6, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x7F, 0xEC, 0x8D, 0x1C, 0x5F, 0x00, 0x00, 0x00, 0xAF, 0x00, 0x00, 0x00, 0x38, 0x00, 0x00, 0x00,
		0x0B, 0x00, 0x00, 0x00, 0x3C, 0x00, 0x3C, 0x00, 0x09, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x5B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x1B, 0x2E, 0x8E, 0x1C, 0x12, 0x03, 0x00, 0x00, 0xFC, 0x03, 0x00, 0x00, 0x43, 0x00, 0x00, 0x00,
		0x2E, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x2E, 0x00, 0x52, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xB3, 0xE6, 0x8E, 0x1C, 0x71, 0x07, 0x00, 0x00, 0x63, 0x05, 0x00, 0x00, 0x71, 0x00, 0x00, 0x00,
		0x63, 0x00, 0x00, 0x00, 0x93, 0x00, 0x1A, 0x00, 0x71, 0x00, 0x9F, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x1A, 0xF5, 0x8E, 0x1C, 0x6F, 0x04, 0x00, 0x00, 0x9E, 0x03, 0x00, 0x00, 0xD4, 0x00, 0x00, 0x00,
		0x13, 0x00, 0x00, 0x00, 0xB6, 0x00, 0x1D, 0x00, 0x43, 0x00, 0x03, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xCC, 0xEB, 0x8F, 0x1C, 0x2E, 0x0A, 0x00, 0x00, 0xD9, 0x30, 0x00, 0x00, 0xE7, 0x00, 0x00, 0x00,
		0xCB, 0x01, 0x00, 0x00, 0x4F, 0x01, 0xAD, 0x00, 0xCE, 0x02, 0x25, 0x00, 0x1D, 0x00, 0xB2, 0x93,
		0x00, 0x00, 0x00, 0x00, 0x85, 0x00, 0x00, 0x00, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	check uint32 = 0x474C4335
)

func TestCRC(t *testing.T) {
	h := NewHash()
	h.Write(data)

	if s := h.Sum32(); s != check {
		t.Fatalf("expected checksum %08x, got %08x", check, s)
	}
}
