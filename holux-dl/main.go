package main

import (
	"fmt"
	"holux"
	"log"
)

func main() {
	c, err := holux.Connect()

	if err != nil {
		log.Println(err)
		return
	}
	c.Hello()
	defer c.Bye()

	index, err := c.GetIndex()
	if err != nil {
		log.Fatalf("Got error %v, arborting", err)
	}

	for _, row := range index {
		fmt.Println(row)
	}
}
