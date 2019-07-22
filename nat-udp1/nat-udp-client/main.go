package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
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

func client(ctx context.Context, address string) (err error) {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return
	}

	localHost := conn.LocalAddr()
	publicHost := stunnedPublicIP()

	log.Printf("Local UDP client address : %s \n", localHost)
	log.Printf("Public UDP client address : %s \n", publicHost)

	defer conn.Close()
	buffer := make([]byte, 1024)

	RemoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	_, err = conn.WriteToUDP([]byte("Hello UDP peer!"), RemoteAddr)
	n, addr, err := conn.ReadFromUDP(buffer)
	fmt.Println("UDP Peer : ", addr)
	fmt.Println("Received from UDP peer : ", string(buffer[:n]))

	return err
}

/*
func client(ctx context.Context, address string) (err error) {
	RemoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", nil)

	//	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		return
	}

	log.Printf("Established connection to %s \n", address)
	log.Printf("remote %s \n", RemoteAddr)

	//log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	localHost := conn.LocalAddr()
	publicHost := stunnedPublicIP()

	log.Printf("Local UDP client address : %s \n", localHost)
	log.Printf("Public UDP client address : %s \n", publicHost)

	defer conn.Close()

	// write a message to server
	message := []byte("Hello UDP server!")

	_, err = conn.WriteToUDP(message, RemoteAddr)

	if err != nil {
		log.Println(err)
	}

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP Server : ", addr)
	fmt.Println("Received from UDP server : ", string(buffer[:n]))

	peer, err := net.ResolveUDPAddr("udp", string(buffer[:n]))
	if err != nil {
		return
	}
	_, err = conn.WriteToUDP([]byte("Hello UDP peer!"), peer)
	buffer = make([]byte, 1024)
	n, addr, err = conn.ReadFromUDP(buffer)
	fmt.Println("UDP Peer : ", addr)
	fmt.Println("Received from UDP peer : ", string(buffer[:n]))

	return err
}
*/
var (
	hostname = flag.String("host", "0.0.0.0", "host to listen to")
	portnum  = flag.String("port", "6000", "port to listen to")
)

func main() {
	flag.Parse()
	hostName := *hostname
	portNum := *portnum

	address := hostName + ":" + portNum
	ctx := context.Background()
	err := client(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
}
