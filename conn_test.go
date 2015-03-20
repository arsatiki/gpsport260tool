package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestChecksum(t *testing.T) {
	msg := []byte("PHLX810")
	if checksum(msg) != 0x35 {
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
		t.Fatal("error during splitting")
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
