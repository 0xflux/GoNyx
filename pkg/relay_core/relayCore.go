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
	_, err := getLocalBinding(global.RelayPort)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print()
	}
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
