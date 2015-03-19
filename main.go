package main

import (
	"bufio"
	"fmt"
	"github.com/arsatiki/term"
	"io"
	"time"
)

const NAP = 150 * time.Millisecond

func checksum(cmd []byte) byte {
	var c byte

	for _, b := range cmd {
		c ^= b
	}
	return c
}

func send(t io.Writer, cmd string) {
	f := "$%s*%02X\r\n"
	cs := checksum([]byte(cmd))
	fmt.Fprintf(t, f, cmd, cs)
}

func receive(s *bufio.Scanner) []byte {
	s.Scan()
	buf := s.Bytes()
	return buf[1 : len(buf)-3]
}

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

func main() {
	// TODO: Stop bits set to 2 for sum reason
	t, err := term.Open("/dev/cu.usbserial",
		term.Speed(38400), term.RawMode)
	if err != nil {
		// TODO LOG
		fmt.Println(err)
		return
	}
	defer t.Close()

	t.SetRTS(true)
	t.SetDTR(true)

	scanner := bufio.NewScanner(t)

	hello(t, scanner)
	defer bye(t, scanner)

	send(t, "PHLX832")
	fmt.Printf("%s\n", receive(scanner))
	
}


