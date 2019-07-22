package main

import (
	"fmt"

	"github.com/ccding/go-stun/stun"
)

func main() {
	s := stun.NewClient()
	s.SetServerAddr("stun.stunprotocol.org:3478")
	//	s.SetVerbose(true)

	for i := 0; i < 99; i++ {
		nat, host, _ := s.Discover()
		fmt.Println(nat, host)
	}
}
