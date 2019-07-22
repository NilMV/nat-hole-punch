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

func client(hostname string, portnum int) {
	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:7070")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	RemoteEP := net.UDPAddr{IP: net.ParseIP(hostname), Port: portnum}
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", localAddr, &RemoteEP)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("Local IP:", conn.LocalAddr())
	log.Println("Remote IP:", conn.RemoteAddr())

	log.Println("Public IP:", stunnedPublicIP())

	defer conn.Close()
	recvBuf := make([]byte, 1024)

	for i := 0; i < 99; i++ {
		n, err := conn.Write([]byte("Ping"))
		if err != nil {
			log.Println(err)
			continue
		}
		if len(recvBuf) != n {
			log.Printf("data size is %d, but sent data size is %d", len(recvBuf), n)
		}
		log.Println("sent")
		n, err = conn.Read(recvBuf)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Received data: %s", string(recvBuf[:n]))
	}
}

func main() {
	flag.Parse()
	hostName := *hostname
	portNum := *portnum

	//address := hostName + ":" + strconv.Itoa(portNum)

	client(hostName, portNum)
}
