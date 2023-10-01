package connectionHandlers

import (
	"fmt"
	"io"
	"log"
	"net"
)

/*

A library to handle connections between relays and other moving parts of the network.

*/

func ReadHTTPRequest(conn net.Conn) ([]byte, error) {
	buff := make([]byte, 1024)
	var data []byte // to allow for any overflows of data in

	// keep reading 1024 bytes until we reach EOF, in which case break. Append all bytes to data variable
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return nil, err
		}

		data = append(data, buff[:n]...)
	}

	return data, nil

}

func SendConnectionToRelay(msg []byte, ip string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = conn.Write(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
