package main

/*
import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var clientStorage map[string]clientEntry

func storeClient(clientDat clientEntry) {
	clientStorage[clientDat.ClientSes] = clientDat
}

func matchClient(conn *net.UDPConn, client1 clientEntry, client2 clientEntry) {
	message := []byte(client2.ClientSes)
	client1Local, err := net.ResolveUDPAddr("udp", client1.LocalAddress)
	if err != nil {
		log.Fatal(err)
	}
	client1Public, err := net.ResolveUDPAddr("udp", client1.PublicAddress)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteToUDP(message, client1Local)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.WriteToUDP(message, client1Public)
	if err != nil {
		log.Println(err)
	}

*/ /*


	message = []byte(client.String())
	_, err = conn.WriteToUDP(message, friendlyPeer)
	if err != nil {
		log.Println(err)
	}*/
/*}

func checkForFriend(conn *net.UDPConn, client clientEntry) {
	if client2, ok := clientStorage[client.ClientSes]; ok {
		matchClient(conn, client, client2)
	}
	storeClient(client)
}

func handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err)
	}

	var clientDat clientEntry

	if err := json.Unmarshal(buffer, &clientDat); err != nil {
		panic(err)
	}
	fmt.Println(clientDat)
	checkForFriend(conn, clientDat)

	if err != nil {
		log.Println(err)
	}
}
*/
