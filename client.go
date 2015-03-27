package holux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/arsatiki/term"
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
	c.ReadReply("PHLX852") // ReadReplyString?

	c.Send("PHLX826")
	c.ReadReply("PHLX859")

	// TODO: Whatabout 832, 861 (firmware)
	// TODO: Handle errors before going here.
	c.SetHighSpeed(921600)
	time.Sleep(NAP)
}

func (c Client) Bye() {
	c.Send("PHLX827")
	c.ReadReply("PHLX860")
	// Slow down port?
}

func (c Client) GetIndex() ([]Index, error) {
	var count int64

	c.Send("PHLX701")
	err := c.ReadInt("PHLX601", &count)
	c.Sendf("PHLX702,0,%d", count)
	err = c.ReadReply("PHLX900,702,3")

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

func (c Client) GetTrack(offset, count int64) (Track, error) {
	var start int64

	c.Sendf("PHLX703,%d,%d", offset, count)
	err := c.ReadReply("PHLX900,703,3")

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
		c.ackBlockHeader()
		block := file[offset : offset+sz]
		// TODO Separate conn into an actual object?
		err = c.ReadBlock(block, bcrc)

		if err != nil {
			// TODO REQUEST RESEND
			panic("Argh")
		} else {
			remaining -= sz
		}
	}
	if cs := CRCChecksum(file); cs != fcrc {
		err = fmt.Errorf("expected file checksum %08x, got %08x", fcrc, cs)
	}
	return file, err
}

func (c Client) readFileHeader() (int64, uint32, error) {
	// expect 901, size, crc
	return 0, 0, nil
}

func (c Client) ackFileHeader() {
	// write 900,901,3
}

func (c Client) readBlockHeader() (int64, int64, uint32, error) {
	// expect 902, offset, len, crc
	return 0, 0, 0, nil
}

func (c Client) ackBlockHeader() {
	// write 900,902,3
}

func (c Client) ackBlock() {
	// write 900,902,3
}
