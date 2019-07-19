package main

import (
	"flag"
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

var (
	hostname = flag.String("host", "0.0.0.0", "host to listen to")
	portnum  = flag.String("port", "6000", "port to listen to")
)

func main() {
	flag.Parse()

	hostName := *hostname
	portNum := *portnum
	service := hostName + ":" + portNum
	clientStorage = make(map[string]*net.UDPAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", service)

	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer ln.Close()

	for {
		handleUDPConnection(ln)
	}

}
