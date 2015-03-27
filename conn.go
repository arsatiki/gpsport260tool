package holux

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const CMDFMT = "$%s*%02X\r\n"

// Conn manages the line-based and binary transmission with any
// io.Reader. In particular, that Reader can be a serial connection
// with the GPS tracker.
type Conn struct {
	rw    io.ReadWriter
	lines *bufio.Scanner
}

func NewConn(rw io.ReadWriter) Conn {
	lines := bufio.NewScanner(rw)
	lines.Split(split)
	return Conn{rw, lines}
}

// TODO: These two should take errors.
func (c Conn) Send(cmd string) {
	cs := foldXOR([]byte(cmd))
	fmt.Fprintf(c.rw, CMDFMT, cmd, cs)
}

func (c Conn) Sendf(format string, a ...interface{}) {
	c.Send(fmt.Sprintf(format, a...))
}

// Receive checks that message from the device matches an expected prefix
func (c Conn) Receive(p string) error {
	c.lines.Scan()
	err := c.lines.Err()
	s := c.lines.Text()

	if err == nil && !strings.HasPrefix(s, p) {
		err = fmt.Errorf("prefix %s not found in %s", p, s)
	}
	return err
}

// Receivef parses the message from the device according to the format
// string.
func (c Conn) Receivef(format string, a ...interface{}) error {
	c.lines.Scan()
	err := c.lines.Err()

	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(c.lines.Text(), format, a...)
	return err
}

// TODO: All the information about n is carried by len(block)
func (c Conn) ReadBlock(block []byte, checksum uint32) error {
	n := int64(len(block))
	h := NewHash()
	src := io.TeeReader(c.rw, h)
	dst := bytes.NewBuffer(block)

	_, err := io.CopyN(dst, src, n)

	if cs := h.Sum32(); err == nil && cs != checksum {
		err = fmt.Errorf("expected block CRC %08x, got %08x", checksum, cs)
	}
	return err
}

func foldXOR(cmd []byte) byte {
	var c byte

	for _, b := range cmd {
		c ^= b
	}
	return c
}

// split validates the checksum and strips out
// $ and other bits from the command
func split(data []byte, atEOF bool) (int, []byte, error) {
	advance, token, err := bufio.ScanLines(data, atEOF)
	if err != nil || token == nil {
		return advance, token, err
	}

	L := len(token)

	if token[0] != '$' || token[L-3] != '*' || L < 3 {
		return advance, token, fmt.Errorf(
			"format mismatch in command %s",
			token)
	}

	cmd := token[1 : L-3]
	check := string(token[L-2 : L])
	cs, err := strconv.ParseUint(check, 16, 8)

	if err != nil {
		return advance, token, fmt.Errorf("bad checksum: %s", check)
	}

	if foldXOR(cmd) != byte(cs) {
		return advance, cmd, fmt.Errorf(
			"checksum mismatch: computed %02x, received %02x",
			foldXOR(cmd),
			cs)
	}

	return advance, cmd, err
}
