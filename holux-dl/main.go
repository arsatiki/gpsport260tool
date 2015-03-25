package main

import (
	"fmt"
	"holux"
)

func main() {
	c, err := holux.Connect()

	if err != nil {
		// TODO LOG
		fmt.Println(err)
		return
	}
	c.Hello()
	defer c.Bye()

	c.Send("PHLX701")
	tcount, _ := c.ReadReply1("PHLX601")

	c.Send(fmt.Sprintf("PHLX702,0,%d", tcount))
	c.ReadReply("PHLX900,702,3")
	ilen, icsum, _ := c.ReadReply1H("PHLX901") // index size, checksum

	c.Send("PHLX900,901,3")
	_, _, _, _ = c.ReadReply3("PHLX902")

	c.Send("PHLX900,902,3")
	index, err := c.ReadIndex(ilen, icsum)

	for _, row := range index {
		fmt.Println(row)
	}
}
