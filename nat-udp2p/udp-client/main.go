package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gortc/stun"
)

func stunnedPublicIP() string {
	var publicHost string
	// Creating a "connection" to STUN server.
	c, err := stun.Dial("udp", "stun.l.google.com:19302")
	if err != nil {
		panic(err)
	}
	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN server, waiting for response message.
	if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			panic(res.Error)
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			panic(err)
		}
		fmt.Println("your IP is", xorAddr.IP, xorAddr.Port)
		publicHost = xorAddr.IP.String() + ":" + strconv.Itoa(xorAddr.Port)
	}); err != nil {
		panic(err)
	}

	return publicHost
}

var (
	hostname = flag.String("host", "0.0.0.0", "host to listen to")
	portnum  = flag.Int("port", 6000, "port to listen to")
)

func client(address string) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("Local IP:", conn.LocalAddr())
	log.Println("Public IP:", stunnedPublicIP())

	defer conn.Close()

	n, err := conn.Write([]byte("Ping"))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	buf := make([]byte, 1024)
	if len(buf) != n {
		log.Printf("data size is %d, but sent data size is %d", len(buf), n)
	}

	recvBuf := make([]byte, 1024)

	n, err = conn.Read(recvBuf)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Printf("Received data: %s", string(recvBuf[:n]))
}

func main() {
	flag.Parse()
	hostName := *hostname
	portNum := *portnum

	address := hostName + ":" + strconv.Itoa(portNum)

	client(address)
}
