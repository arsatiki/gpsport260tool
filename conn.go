package main

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/arsatiki/term"
)

const (
	NAP    = 150 * time.Millisecond
	CMDFMT = "$%7s*%02X"
)

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
	// TODO SPLITFUNC with checksum audit
	return &Conn{t: t, lines: lines}, nil
}

func (c Conn) Close() {
	c.t.Close()
}

// TODO: These two should take errors.
func (c Conn) Send(cmd string) {
	cs := checksum([]byte(cmd))
	fmt.Fprintf(c.t, CMDFMT, cmd, cs)
	c.t.Write([]byte("\r\n"))
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

/*func (c Conn) ReadReplyArg(expected string) (int, error) {
	comma := []byte(",")
	reply := c.receive()
	parts := bytes.SplitN(reply, comma, 1)

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
	var (
		cmd []byte
		cs  byte
	)

	advance, token, err := bufio.ScanLines(data, atEOF)
	if err != nil || token == nil {
		return advance, token, err
	}

	// This restringifying seems silly
	_, err = fmt.Sscanf(string(token), CMDFMT, &cmd, &cs)
	if err != nil {
		return advance, token, err
	}
	if checksum(cmd) != cs {
		return advance, cmd, fmt.Errorf(
			"checksum mismatch: computed %02x, received %02x",
			checksum(cmd),
			cs)
	}

	return advance, cmd, err
}
