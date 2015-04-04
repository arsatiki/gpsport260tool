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

	for k, row := range index {
		fmt.Println(row)
		track, err := c.GetTrack(row.Offset, row.Size)
		if err != nil {
			log.Fatal("Got error %v while reading track %d", err, k)
		}
		for _, t := range track {
			fmt.Println(t)
		}
	}
}

func init() {
	// TODO: Prepare SQL queries.
}