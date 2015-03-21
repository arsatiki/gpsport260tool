package holux

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/arsatiki/term"
)

const (
	NAP    = 150 * time.Millisecond
	CMDFMT = "$%s*%02X\r\n"
)

// Conn manages the line-based and binary transmission with the GPS tracker.
// TODO: Conn should use a generic reader for easier testing, rite?
type Conn struct {
	t     *term.Term
	lines *bufio.Scanner
}

func Connect() (*Conn, error) {
	t, err := term.Open("/dev/cu.usbserial",
		term.Speed(38400), term.RawMode)

	if err != nil {
		return nil, err
	}

	t.SetRTS(true)
	t.SetDTR(true)

	lines := bufio.NewScanner(t)
	return &Conn{t: t, lines: lines}, nil
}

func (c Conn) Close() {
	c.t.Close()
}

// TODO: These two should take errors.
func (c Conn) Send(cmd string) {
	cs := checksum([]byte(cmd))
	fmt.Fprintf(c.t, CMDFMT, cmd, cs)
}

func (c Conn) receive() []byte {
	// TODO ERROR?
	c.lines.Scan()
	buf := c.lines.Bytes()
	return buf[1 : len(buf)-3]
}

func (c Conn) ReadReply(expected string) error {
	reply := c.receive()
	if !bytes.Equal(reply, []byte(expected)) {
		return fmt.Errorf(
			"bad reply, expected %s, got %s",
			expected,
			reply)
	}

	return nil
}

/*
func (c Conn) ReadReplyArg(expected string) (int, error) {

	reply := string(c.receive())
	comma := []byte(",")
	reply := c.receive()
	parts := bytes.SplitN(reply, comma, 1)
	// TODO: parts[0]
}*/

/*
func hello(t *term.Term, s *bufio.Scanner) {
	send(t, "PHLX810")
	fmt.Printf("%s\n", receive(s))

	send(t, "PHLX826")
	fmt.Printf("%s\n", receive(s))
	t.SetHighSpeed(921600)
	time.Sleep(NAP)
}

func bye(t *term.Term, s *bufio.Scanner) {
	send(t, "PHLX827")
	fmt.Printf("%s\n", receive(s))
}
*/

func checksum(cmd []byte) byte {
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

	if checksum(cmd) != byte(cs) {
		return advance, cmd, fmt.Errorf(
			"checksum mismatch: computed %02x, received %02x",
			checksum(cmd),
			cs)
	}

	return advance, cmd, err
}
