package main

import "github.com/arsatiki/term"

func main() {
	t := term.Open("/dev/cu.usbserial", term.Speed(38400))
	t.setRTS(true)
	t.setDTR(true)
}