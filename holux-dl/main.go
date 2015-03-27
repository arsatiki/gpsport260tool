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

	index, err := c.GetIndex()
	if err != nil {
		fmt.Printf("Got error %v, arborting", err)
	}
	
	for _, row := range index {
		fmt.Println(row)
	}
}
