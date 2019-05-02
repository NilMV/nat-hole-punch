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

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		return
	}

	log.Printf("Established connection to %s \n", address)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	// write a message to server
	message := []byte("Hello UDP server!")

	_, err = conn.Write(message)

	if err != nil {
		log.Println(err)
	}

	// receive message from server
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP Server : ", addr)
	fmt.Println("Received from UDP server : ", string(buffer[:n]))
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
