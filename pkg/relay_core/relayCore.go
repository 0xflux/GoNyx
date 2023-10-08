package relay_core

import (
	cryptolocal "GoNyx/pkg/crypto"
	"GoNyx/pkg/global"
	"bufio"
	"crypto/ecdh"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/*
Nyx Relay core engine.
*/

// StartRelay entry point to starting the relay
func StartRelay() {
	// instantiate the relay
	this := NewRelay()
	fmt.Println("Server public fingerprint: ", this.PublicKeyHash)

	// start the listeners for the relay
	manageListeners(this)
}

func manageListeners(this *Relay) {
	// to handle concurrency and prevent main from exiting
	stop := make(chan os.Signal, 1)

	// listen to interrupt signals
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// start the various listeners on the relay
	go func() {
		startListener("relay", this) // start listener for relay
	}()
	go func() {
		startListener("negotiation", this) // for crypto and route negotiation
	}()

	// wait for signal interrupt
	<-stop
	fmt.Println("Stopping..")
}

// startListener to be used as a goroutine to listen on a certain binding
func startListener(t string, this *Relay) {
	var listener net.Listener
	switch t {
	case "relay":
		listener, _ = getLocalBinding(global.RelayPort)

	case "negotiation":
		listener, _ = net.Listen("tcp", global.RelayCryptoNegotiation)

	default:
		log.Fatal("Requires argument in function")
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
		go handleConnection(conn, this)
	}
}

// handleConnection will handle inbound http requests
func handleConnection(conn net.Conn, this *Relay) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	//buff := make([]byte, 1024) // what happens if this overflows? Err?
	//n, err := conn.Read(buff)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//fmt.Println("Data received: ", string(buff[:n]))

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Println("Error reading HTTP request:", err)
		return
	}

	if req.Body == nil {
		log.Println("No body in request")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body")
		}
	}(req.Body)

	// Read the body (which should contain your JSON)
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}

	var msg Relay
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	// Now you can access fields from msg, e.g., msg.PrivateKey
	fmt.Println("Received Public Key:", msg.PublicKey)

	fmt.Println("Calculating secret....")
	pub, err := ecdh.P521().NewPublicKey(msg.PublicKey)
	if err != nil {
		log.Println("Error calculating public key")
	}

	priv, err := ecdh.P521().NewPrivateKey(this.PrivateKey)
	if err != nil {
		log.Println("Error calculating private key")
	}

	res, err := cryptolocal.ComputeSharedSecret(pub, priv)
	if err != nil {
		fmt.Println("Error after func: ", err)
	}
	fmt.Println("Shared secret generated, result: ", res)
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
