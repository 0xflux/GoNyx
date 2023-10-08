package clientcore

import (
	connectionHandlers "GoNyx/pkg/connection_handlers"
	cryptolocal "GoNyx/pkg/crypto"
	"GoNyx/pkg/global"
	"crypto/ecdh"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

/*
Nyx Client core engine.
*/

// StartClient start the engine on the client
func StartClient() {

	/*
		TODO: probably want a concurrent func here handling routing before it even starts
		listening, this will then maintain live circuits for the client to use.
	*/

	go secretTest()

	// start listener
	startListener(global.ClientListenAddr)
}

// startListener to be used as a goroutine to listen on a certain binding
func startListener(addr string) {
	listener, err := net.Listen("tcp", addr)
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

// secretTest tests the secret sharing algorithm
func secretTest() {
	privateKey, _, err := cryptolocal.NewECDHKeyPair()
	if err != nil {
		log.Println("Error creating cryptographic keys for communication. ", err)
	}

	log.Println("Calculating secret client side.")

	// proof of concept by using the current public key of the server
	pubKey := []byte{4, 1, 24, 212, 73, 208, 70, 152, 8, 253, 146, 236, 224, 0, 179, 167, 247, 187, 33, 193, 207, 210,
		161, 156, 221, 89, 147, 167, 112, 128, 28, 207, 41, 53, 74, 97, 68, 204, 154, 69, 209, 23, 40, 162, 40, 228, 49,
		241, 90, 13, 53, 34, 210, 45, 222, 174, 208, 145, 233, 200, 51, 243, 113, 123, 136, 100, 238, 0, 46, 224, 213,
		80, 84, 79, 107, 46, 232, 207, 88, 212, 230, 221, 107, 193, 195, 95, 140, 89, 65, 51, 228, 250, 217, 19, 89,
		253, 21, 104, 88, 200, 98, 115, 172, 132, 4, 215, 229, 164, 207, 239, 92, 41, 101, 96, 164, 190, 81, 175, 238,
		199, 145, 154, 123, 107, 220, 179, 230, 140, 140, 226, 73, 98, 184}

	p, err := ecdh.P521().NewPublicKey(pubKey)
	if err != nil {
		log.Fatal("Error making public key - ", err)
	}

	if secret, err := privateKey.ECDH(p); err != nil {
		log.Fatalf("Error finding secret, %v\n", err)
	} else {
		fmt.Println("Secret is: ", secret)
		if cipher, err := cryptolocal.EncryptCommunication(secret, []byte("hello world")); err != nil {
			fmt.Println("Error in encrypting comms and sending. ", err)
		} else {
			fmt.Println("Successfully encrypted message: %v", cipher)
			// decrypt
			if plain, err := cryptolocal.DecryptCommunication(cipher, secret); err != nil {
				log.Fatal("Error decrypting, quitting.")
			} else {
				fmt.Println("Plaintext: ", plain)
			}
		}

	}

	//url := fmt.Sprintf("http://%s:%v", global.ListenIP, global.NegotiationPort)
	//msg := &relay_core.Relay{PrivateKey: privateKey.Bytes(), PublicKey: publicKey.Bytes()}
	//jsonData, err := json.Marshal(msg)
	//if err != nil {
	//	log.Printf("Error marshal json, %s\n", err)
	//}
	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	//if err != nil {
	//	log.Println("Error sending public key:", err)
	//	return
	//}
	//defer resp.Body.Close()

	// log.Println("Public Key sent with status:", resp.Status)
}

// process connection to proxy
func processConnection(conn net.Conn) {

	// ensures we close the connection in the right scope.
	// therefore, cannot use the connection outside of this function; can use within
	// nested functions however.
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	// complete SOCKS5 handshake
	targetAddress, err := outboundSocksHandshake(conn)
	if err != nil {
		log.Printf("Handshake error: %v. targetAddress: %s.\nAborting request.\n", targetAddress, err)
		return
	}

	addrType := global.ClassifyAddress(targetAddress)
	switch addrType {
	case "domain":
		// handle domain connection
		handleDomainConnection(targetAddress, conn)
	case "IPv4":
		// (route via exits)
	case "IPv6":
		log.Println("IPv6 not supported. Aborting connection.")
		return
	}

}

// handle domain connection after SOCKS handshake
func handleDomainConnection(targetAddress string, conn net.Conn) {

	if strings.HasSuffix(strings.Split(targetAddress, ":")[0], ".nyx") {
		fmt.Println("Nyx address")
		// handle .nyx protocol
	} else {
		// handle clear web connection routing
		// exit will use net.Dial() not the client I think
		httpData, err := connectionHandlers.ReadHTTPRequest(conn)
		if err != nil {
			// currently generating lots of errors, I think it's expected? Silencing for now...
			// log.Println(err)
			// possibly producing errors because its TLS. TODO
			return
		}

		// test case for debugging, hxxp://something[.]com
		if strings.Contains(httpData.Host, "something") {
			go connectionHandlers.SendConnectionToRelay(httpData, global.ListenIP, global.RelayPort)
		}
	}

}

// handle SOCKS5 handshake
func outboundSocksHandshake(conn net.Conn) (string, error) {
	/*
		Good info on this:
		https://medium.com/@nimit95/socks-5-a-proxy-protocol-b741d3bec66c

		IMPORTANT:
		Once a byte is read from a net.Conn object, it is consumed and subsequent reads will
		read the next bytes. Kinda like reading from a file: as you read, your position increments.
	*/

	// first two bytes, SOCKS version & authentication method
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		log.Printf("Error consuming first 2 bytes of connection data. %v\n", err)
		return "", err
	}

	// check socks is running as v5
	if buf[0] != 5 {
		log.Printf("Wrong SOCKS version detected, make sure you use SOCKSv5. SOCKS detected is version %v. "+
			"Refusing connection.\n", buf[0])
		return "", errors.New("invalid SOCKS version")
	}

	// number of authentication methods the client supports
	numAuthMethods := int(buf[1])
	_, err = io.ReadFull(conn, buf[:numAuthMethods])
	if err != nil {
		log.Printf("Error consuming methods %v.", err)
		return "", err
	}

	// the handshake requires a response at this stage,
	// the second byte being the authentication method by the proxy
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
		log.Printf("Error reading request data, %v\n", err)
		return "", err
	}

	// only support TCP/IP stream connection for now
	if buf[1] != 1 {
		log.Printf("Unsupported command, expecting 0x1, found: %v.\n", buf[1])
		return "", errors.New("unsupported command, expecting 0x1")
	}

	var targetAddress string
	switch buf[3] {
	case 1: //ipv4
		_, err := io.ReadFull(conn, buf[:6])
		if err != nil {
			log.Printf("Error reading IPv4: %v\n", err)
			return "", err
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
			log.Printf("Error reading domain length: %v\n", err)
			return "", err
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
			log.Printf("Error reading domain: %v\n", err)
			return "", err
		}

		// fmt.Printf("\nbuf as hex %02x, buf as str: %s\n", buf, buf)

		// decode the domain and port from the buffer, the port is found between the end of the
		// domain + 2 bytes (i.e. the next 2 bytes after the end of the domain)
		targetAddress = fmt.Sprintf("%s:%d", buf[:domainLength], binary.BigEndian.Uint16(buf[domainLength:domainLength+2]))

	case 4: // ipv6 refuse
		log.Println("IPv6 addresses not supported.")
		return "", errors.New("ipv6 addresses not supported")
	default:
		log.Println("Invalid address type.")
		return "", errors.New("invalid address type")
	}

	// send success response to client
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return "", err
	}

	if targetAddress != "" {
		return targetAddress, nil
	}

	// if somehow target address is empty string
	return targetAddress, errors.New("invalid target address, debug")
}
