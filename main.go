package main

import (
	"fmt"
)

func main() {
	c, err := Connect()

	if err != nil {
		// TODO LOG
		fmt.Println(err)
		return
	}
	defer c.Close()

	/*
		hello(t, scanner)
		defer bye(t, scanner)

		send(t, "PHLX832")
		receive(scanner)

		// List trakcs
		send(t, "PHLX701")
		reply := receive(scanner)
		cmd, count := bytes.Split(reply, []byte(","))
	*/
}
