package clientcore

import (
	"GoNyx/pkg/global"
	"fmt"
	"net"
)

/*
Nyx Client core engine.
*/

// start the engine on the client
func StartClient() {
	listener, err := net.Listen("tcp", global.ClientListenAddr)
	if err != nil {
		fmt.Printf("Error listening on tcp at address: %v. %v\n", global.ClientListenAddr, err)
	}
	defer listener.Close()

	for {
		browserConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error listening to browser.%v\n", err)
		}

		// TODO: pick these outputs apart, obviously pointing to some structure in memory
		fmt.Printf("Listener: %v\n", listener)
		fmt.Printf("Browser conn: %v\n", browserConn)

	}
}

/*
initiate circuit
*/
