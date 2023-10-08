package connectionHandlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

/*

A library to handle connections between relays and other moving parts of the network.

*/

// ReadHTTPRequest read a HTTP request from a connection stream
func ReadHTTPRequest(conn net.Conn) (*http.Request, error) {
	// note, data is nil for TLS traffic, as cannot read the data.
	data, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SendEncryptedConnectionToRelay will send an encrypted envelope via TCP to the relay, this will be the normal routing
// for traffic through the Nyx Net. However, if data is leaving a relay as a HTTP request, instead use function
// SendHTTPRequest, which forwards the final message as a http request, not a raw TCP request as is with this func.
func SendEncryptedConnectionToRelay(payload []byte, ip string, port int) {
	url := fmt.Sprintf("%s:%v", ip, port)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Println("Error dialing url to send encrypted message. ", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection from sending http data. ", err)
			return
		}
	}(conn)

	// send over the connection
	_, err = conn.Write(payload)
	if err != nil {
		log.Println("Error writing connection data. ", err)
		return
	}

}

// SendHTTPRequest sends a HTTP request from the machine to its destination. In most cases, this should be out of the
// final relay to the internet, if it is not bound for a .nyx domain. I think... For relay to relay communication, use
// SendEncryptedConnectionToRelay instead.
func SendHTTPRequest(msg *http.Request, ip string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println("Error connecting to relay, is it online? Err: ", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error cleaning up connection stream: ", err)
		}
	}()

	// convert the request to its byte form
	data, err := httputil.DumpRequest(msg, true)
	if err != nil {
		log.Println("Error converting request to bytes: ", err)
		return
	}
	_, err = conn.Write(data)

	if err != nil {
		log.Println("Error writing data to connection stream ", err)
		return
	}

	fmt.Println("Connection data sent")
}
