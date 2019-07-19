package main

import (
	"fmt"

	"github.com/ccding/go-stun/stun"
)

func main() {
	s := stun.NewClient()
	s.SetServerAddr("stun.l.google.com:19302")
	s.SetVerbose(true)
	nat, host, _ := s.Discover()
	fmt.Print(nat, host)
}
