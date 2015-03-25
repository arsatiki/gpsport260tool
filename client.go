package holux

import (
	"github.com/arsatiki/term"
	"time"
)

const NAP = 150 * time.Millisecond

// Client describes the protocol level interactions with
// the tracker.
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

	c := &Client{NewConn(t), t}
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
