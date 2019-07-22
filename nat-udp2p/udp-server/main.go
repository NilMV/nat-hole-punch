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

func server(address *net.UDPAddr) {
	updLn, err := net.ListenUDP("udp", address)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	log.Println("Starting udp server...")

	log.Println("Local IP:", updLn.LocalAddr())
	log.Println("Public IP:", stunnedPublicIP())

	for {
		n, addr, err := updLn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		go func() {
			log.Printf("Reciving data: %s from %s", string(buf[:n]), addr.String())

			log.Printf("Sending data..")
			updLn.WriteTo([]byte("Pong"), addr)
			log.Printf("Complete Sending data..")
		}()
	}
}

func main() {
	flag.Parse()
	hostName := *hostname
	portNum := *portnum

	//	address := hostName + ":" + portNum
	//ctx := context.Background()
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(hostName),
		Port: portNum,
	}
	server(udpAddr)
	/*
		err := client(ctx, address)
		if err != nil {
			log.Fatal(err)
		}
	*/
}
