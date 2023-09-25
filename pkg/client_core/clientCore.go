package clientcore

import (
	"GoNyx/pkg/global"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

/*
Nyx Client core engine.
*/

// start the engine on the client
func StartClient() {

	/*
		TODO: probably want a concurent func here handling routing before it even starts
		listening, this will then maintain live circuits for the client to use.
	*/

	// start listener
	listener, err := net.Listen("tcp", global.ClientListenAddr)
	if err != nil {
		log.Fatalf("Error starting listener, %v\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection, %v\n", err)
			continue
		}

		go processConnection(conn)
	}
}

// process connection to proxy
func processConnection(conn net.Conn) {

	/*
		Good info on this:
		https://medium.com/@nimit95/socks-5-a-proxy-protocol-b741d3bec66c

		IMPORTANT:
		Once a byte is read from a net.Conn object, it is consumed and subsequent reads will
		read the next bytes. Kinda like reading from a file: as you read, your position increments.
	*/
	defer conn.Close()

	// first two bytes, SOCKS version & authentication method
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		fmt.Printf("Error consuming first 2 bytes of connection data. %v\n", err)
		return
	}

	// check socks is running as v5
	if buf[0] != 5 {
		fmt.Printf("Wrong SOCKS version detected, make sure you use SOCKSv5. SOCKS detected is version %v. "+
			"Refusing connection.\n", buf[0])
	}

	// number of authentication methods the client supports
	numAuthMethods := int(buf[1])
	_, err = io.ReadFull(conn, buf[:numAuthMethods])
	if err != nil {
		fmt.Printf("Error consuming methods %v.", err)
		return
	}

	// the handshake requires a response at this stage,
	// the second byte being the authenticaiton method by the proxy
	conn.Write([]byte{0x05, 0x00})

	// Client sends the request packet (0x05, 0x01, 0x00, 0x03, <B_HOST>, <B_PORT>)
	/*
		The Second Byte 0x01 is for the command code. It is one byte.
			0x01: establish a TCP/IP stream connection

			0x02: establish a TCP/IP port binding

			0x03: associate a UDP port

		The Third Byte 0x00 is a reserved byte. It must be 0x00 and 1 byte.

		The Fourth Byte 0x03 is the address type of desired HOST and 1 byte.
			0x01: IPv4 address, followed by 4 bytes IP

			0x03: Domain name, 1 byte for name length, followed by host name

			0x04: IPv6 address, followed by 16 bytes IP

		The last Byte is port number in a network byte order, 2 bytes
	*/

	buf = make([]byte, 1024) // not sure what the safest size is. maybe need better error logging

	_, err = io.ReadFull(conn, buf[:4])
	if err != nil {
		fmt.Printf("Error reading request data, %v\n", err)
		return
	}

	// only support TCP/IP stream connection for now
	if buf[1] != 1 {
		fmt.Printf("Unsupported command, expecting 0x1, found: %v.\n", buf[1])
	}

	var targetAddress string
	switch buf[3] {
	case 1: //ipv4
		_, err := io.ReadFull(conn, buf[:6])
		if err != nil {
			fmt.Printf("Error reading IPv4: %v\n", err)
			return
		}
		targetAddress = fmt.Sprintf("%d.%d.%d.%d:%d", buf[0], buf[1], buf[2], buf[3], binary.BigEndian.Uint16(buf[4:6]))

	case 3: // domain
		/*
			In the SOCKS5 protocol, if the address type is a domain name (ATYP is 0x03),
			the next byte will indicate the length of the domain name.
			Read the length of the byte and stores it as domainLength.
		*/
		_, err = io.ReadFull(conn, buf[:1])
		if err != nil {
			fmt.Printf("Error reading domain length: %v\n", err)
			return
		}
		domainLength := int(buf[0])

		// for _, b := range buf {
		// 	if b != 0x00 {
		// 		fmt.Printf("%02x ", b)
		// 	}
		// }

		// read the domain and port bytes into the buffer
		_, err = io.ReadFull(conn, buf[:domainLength+2]) // domain + next 2 bytes for port
		if err != nil {
			fmt.Printf("Error reading domain: %v\n", err)
			return
		}

		// fmt.Printf("\nbuf as hex %02x, buf as str: %s\n", buf, buf)

		// decode the domain and port from the buffer, the port is found between the end of the
		// domain + 2 bytes (i.e. the next 2 bytes after the end of the domain)
		targetAddress = fmt.Sprintf("%s:%d", buf[:domainLength], binary.BigEndian.Uint16(buf[domainLength:domainLength+2]))

	case 4: // ipv6 refuse
		fmt.Println("IPv6 addresses not supported.")
		return
	default:
		fmt.Println("Invalid address type")
		return
	}

	fmt.Println(targetAddress)
}
