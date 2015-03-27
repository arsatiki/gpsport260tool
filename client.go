package holux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/arsatiki/term"
	"log"
	"time"
)

const NAP = 150 * time.Millisecond

// Client describes the protocol level interactions with
// the tracker. Methods starting with Get do multiple requests
// to the device. Methods starting with read or ack
// do only single requests.
type Client struct {
	Conn
	*term.Term
}

func Connect() (*Client, error) {
	t, err := term.Open("/dev/cu.usbserial",
		term.Speed(38400), term.RawMode)

	if err != nil {
		return nil, err
	}

	c := &Client{Conn: NewConn(t), Term: t}
	c.SetRTS(true)
	c.SetDTR(true)

	return c, nil
}

func (c Client) Close() {
	c.Close()
}

func (c Client) Hello() {
	c.Send("PHLX810")
	c.Receive("PHLX852") // Receivef?

	c.Send("PHLX826")
	c.Receive("PHLX859")

	// TODO: Whatabout 832, 861 (firmware)
	// TODO: Handle errors before going here.
	c.SetHighSpeed(921600)
	time.Sleep(NAP)
}

func (c Client) Bye() {
	c.Send("PHLX827")
	c.Receive("PHLX860")
	// Slow down port?
}

func (c Client) GetIndex() ([]Index, error) {
	var count int64

	c.Send("PHLX701")
	err := c.Receivef("PHLX601,%d", &count)
	if err != nil {
		log.Println("Received error: ", err)
	}

	c.Sendf("PHLX702,0,%d", count)
	err = c.Receive("PHLX900,702,3")

	f, err := c.GetFile()
	if err != nil {
		return nil, err
	}
	index := make([]Index, count)
	buf := bytes.NewBuffer(f)
	err = binary.Read(buf, binary.LittleEndian, index)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (c Client) GetTrack(offset, count uint32) (Track, error) {
	c.Sendf("PHLX703,%d,%d", offset, count)
	err := c.Receive("PHLX900,703,3")

	f, err := c.GetFile()
	if err != nil {
		return nil, err
	}
	track := make([]Trackpoint, count)
	buf := bytes.NewBuffer(f)
	err = binary.Read(buf, binary.LittleEndian, track)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (c Client) GetFile() ([]byte, error) {
	remaining, fcrc, err := c.readFileHeader()
	c.ackFileHeader()

	file := make([]byte, remaining)
	for remaining > 0 {
		offset, sz, bcrc, err := c.readBlockHeader()
		if err != nil {
			log.Fatal("could not read block header", err)
		}
		c.ackBlockHeader()
		block := file[offset : offset+sz]
		// TODO Separate conn into an actual object?
		err = c.ReadBlock(block, bcrc)
		if err != nil {
			log.Fatal("file transfer error: ", err, fcrc)
			// TODO REQUEST RESEND
		} else {
			remaining -= sz
		}
	}
	if cs := CRCChecksum(file); cs != fcrc {
		err = fmt.Errorf("expected file checksum %08x, got %08x", fcrc, cs)
	}
	return file, err
}

func (c Client) readFileHeader() (sz int64, crc uint32, err error) {
	err = c.Receivef("PHLX901,%d,%x", &sz, &crc)
	return sz, crc, err
}

func (c Client) ackFileHeader() {
	c.Send("PHLX900,901,3")
}

func (c Client) readBlockHeader() (start, sz int64, crc uint32, err error) {
	err = c.Receivef("PHLX902,%d,%d,%x", &start, &sz, &crc)
	return start, sz, crc, err
}

func (c Client) ackBlockHeader() {
	c.Send("PHLX900,902,3")
}

func (c Client) ackBlock() {
	c.Send("PHLX900,902,3")
}
