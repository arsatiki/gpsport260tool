package holux

import (
	"bufio"
	"bytes"
	"testing"
)

func TestChecksum(t *testing.T) {
	msg := []byte("PHLX810")
	if foldXOR(msg) != 0x35 {
		t.Fatal("checksum mismatch")
	}
}

func TestSplitSuccess(t *testing.T) {
	data := []byte("$PHLX810*35\r\n")
	buf := bytes.NewBuffer(data)
	s := bufio.NewScanner(buf)
	s.Split(split)
	s.Scan()

	if s.Err() != nil {
		t.Fatal("error during splitting:", s.Err())
	}
	if s.Text() != "PHLX810" {
		t.Fatal("parsed incorrect token")
	}
}

func TestSplitFailure(t *testing.T) {
	data := []byte("$PHLX810*AA\r\n")
	buf := bytes.NewBuffer(data)
	s := bufio.NewScanner(buf)
	s.Split(split)
	s.Scan()

	if s.Err() == nil {
		t.Fatal("splitting should have failed")
	}
}

func TestTrackReadBlock(t *testing.T) {
	data := []byte{0x7F, 0xEC, 0x8D, 0x1C, 0x65, 0xC8, 0x70, 0x42,
		0xC0, 0x88, 0xC7, 0x41, 0x81, 0x00, 0x39, 0x00,
		0x48, 0x05, 0x00, 0x00, 0x1F, 0x00, 0x6A, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	target := make([]byte, 32)

	c := NewConn(bytes.NewBuffer(data))
	err := c.ReadBlock(target, 0x421fe9dd)
	if err != nil {
		t.Fatalf("failed to read a block: %v", err)
	}
}

