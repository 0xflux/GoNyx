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

func SendConnectionToRelay(msg *http.Request, ip string, port int) {
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
