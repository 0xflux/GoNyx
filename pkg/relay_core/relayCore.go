package relay_core

import (
	"GoNyx/pkg/global"
	"errors"
	"fmt"
	"log"
	"net"
)

/*
Nyx Relay core engine.
*/

// StartRelay entry point to starting the relay
func StartRelay() {
	// in local testing we have 3 predefined ports to use, so assign the relay a port number
	listener, err := getLocalBinding(global.RelayPort)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection. ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Data received: ", string(buff[:n]))
}

// gets local bind address on localhost based off of 3 port numbers for debugging
func getLocalBinding(startPort int) (net.Listener, error) {

	maxPort := global.RelayPort + 3 // define max 4 ports we want to listen to in debug on relays

	for port := startPort; port <= maxPort; port++ {
		listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ListenIP, port))
		if err == nil {
			// no error so return the listener
			fmt.Println("Listening on port ", port)
			return listen, nil
		}
	}

	return nil, errors.New("cannot assign free port")
}
