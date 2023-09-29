package global

import "net"

func ClassifyAddress(address string) string {
	// Split address into host and port
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		// try to treat the address as a host without a port
		host = address
	}

	// parse the host part as IP
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.To4() != nil {
			return "IPv4"
		}
		if ip.To16() != nil {
			return "IPv6"
		}
	}

	// if not an IP, treat it as a domain
	return "domain"
}
