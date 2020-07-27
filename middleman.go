package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
func handleClient(client net.Conn, remote net.Conn) {
	//defer client.Close()
	//chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		//chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		//chDone <- true
	}()
	//<-chDone
}
