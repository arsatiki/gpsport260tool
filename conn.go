package holux

import (
	"bufio"
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
	rw     io.ReadWriter
	lines *bufio.Scanner
}

func NewConn(rw io.ReadWriter) Conn {
	lines := bufio.NewScanner(rw)
	lines.Split(split)
	return Conn{rw, lines}
}

// TODO: These two should take errors.
func (c Conn) Send(cmd string) {
	cs := checksum([]byte(cmd))
	fmt.Fprintf(c.rw, CMDFMT, cmd, cs)
}

// ReadReply reads a message from the tracker and validates that
// the prefix matches.
// TODO: There's no string reply reading yet. Consider renaming to
// Read0, Read1, Read3 and ReadS.
func (c Conn) ReadReply(p string) error {
	_, err := c.receiveLine(p)
	return err
}

// ReadReply1 reads a message from the tracker and returns
// the numeric parameter.
func (c Conn) ReadReply1(p string) (int64, error) {
	args, err := c.receiveLine(p)
	if err != nil {
		return 0, err
	}

	vals, err := parseIntArgs(args, 10)
	return vals[0], err
}

// ReadReply3 reads a message from the tracker containing
// two decimal parameters and a checksum.
func (c Conn) ReadReply3(p string) (int64, int64, int64, error) {
	args, err := c.receiveLine(p)
	if err != nil {
		return 0, 0, 0, err
	}

	vals, err := parseIntArgs(args, 10, 10, 16)
	return vals[0], vals[1], vals[2], err

}

// TODO: ReadBinary

// receiveLine expects a prefix of the expected reply as a string
// and returns split out data as args for further processing
func (c Conn) receiveLine(p string) ([]string, error) {
	c.lines.Scan()

	if c.lines.Err() != nil {
		return nil, c.lines.Err()
	}

	reply := c.lines.Text()
	if !strings.HasPrefix(reply, p) {
		return nil, fmt.Errorf(
			"expected prefix %s, got %s",
			p,
			reply)
	}

	parts := strings.Split(reply, ",")
	return parts[1:], nil
}

// parseIntArgs converts a group of integer strings at once.
// The slice always has the same number of elements as the
// number of bases given.
func parseIntArgs(args []string, bases ...int) ([]int64, error) {
	A := len(args)
	B := len(bases)
	vals := make([]int64, B)

	if A != B {
		err := fmt.Errorf("Expected %d args, got %d", B, A)
		return vals, err
	}

	for k, s := range args {
		v, err := strconv.ParseInt(s, bases[k], 0)
		if err != nil {
			return nil, err
		}
		vals[k] = v
	}
	return vals, nil
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
