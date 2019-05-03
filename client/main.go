package main

import (
	"context"
	"fmt"
	"log"
	"net"
)

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
	//	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	//log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

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

func main() {
	hostName := "localhost"
	portNum := "6000"

	address := hostName + ":" + portNum
	ctx := context.Background()
	err := client(ctx, address)
	if err != nil {
		log.Fatal(err)
	}
}
