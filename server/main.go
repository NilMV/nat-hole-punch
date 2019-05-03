package main

import (
	"fmt"
	"log"
	"net"
)

var clientStorage map[string]*net.UDPAddr

func storeClient(client *net.UDPAddr, idSession string) {
	clientStorage[idSession] = client
}

func matchClient(conn *net.UDPConn, client *net.UDPAddr, friendlyPeer *net.UDPAddr) {
	message := []byte(friendlyPeer.String())
	_, err := conn.WriteToUDP(message, client)
	if err != nil {
		log.Println(err)
	}
	message = []byte(client.String())
	_, err = conn.WriteToUDP(message, friendlyPeer)
	if err != nil {
		log.Println(err)
	}
}

func checkForFriend(conn *net.UDPConn, client *net.UDPAddr, idSession string) {
	if friendlyPeer, ok := clientStorage[idSession]; ok {
		matchClient(conn, client, friendlyPeer)
	}
	storeClient(client, idSession)
}

func handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}

	idSession := string(buffer[:n])
	checkForFriend(conn, addr, idSession)

	if err != nil {
		log.Println(err)
	}

}

func main() {
	hostName := "localhost"
	portNum := "6000"
	service := hostName + ":" + portNum
	clientStorage = make(map[string]*net.UDPAddr)

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}

}
