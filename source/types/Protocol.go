package types

import "strconv"

type Protocol string

const (
	ProtocolDNS        Protocol = "dns"
	ProtocolDNSoverTLS Protocol = "dns-over-tls"
	ProtocolHTTPS      Protocol = "https"
	ProtocolHTTP       Protocol = "http"
	ProtocolICMP       Protocol = "icmp"
	ProtocolSSH        Protocol = "ssh"
	ProtocolSOCKS      Protocol = "socks"
	ProtocolTCP        Protocol = "tcp"
	ProtocolUDP        Protocol = "udp"
	ProtocolWHOIS      Protocol = "whois"
)

func IsProtocol(value string) bool {

	switch Protocol(value) {
	case Any,
		ProtocolDNS,
		ProtocolDNSoverTLS,
		ProtocolHTTPS,
		ProtocolHTTP,
		ProtocolICMP,
		ProtocolSSH,
		ProtocolSOCKS,
		ProtocolTCP,
		ProtocolUDP,
		ProtocolWHOIS:
		return true
	}

	return false

}

func NewProtocol() Protocol {
	return Protocol(Any)
}

func ParseProtocol(value string) *Protocol {

	var result *Protocol = nil
	var protocol Protocol

	if protocol.Parse(value) == true {
		result = &protocol
	}

	return result

}

func ToProtocol(value string) Protocol {

	var protocol Protocol

	protocol.Parse(value)

	return protocol

}

func (protocol *Protocol) Parse(value string) bool {

	if value == "*" || value == "all" || value == Any {

		*protocol = Protocol(Any)

		return true

	} else if value == "https" {

		*protocol = ProtocolHTTPS

		return true

	} else if value == "http" {

		*protocol = ProtocolHTTP

		return true

	} else if value == "dns" {

		*protocol = ProtocolDNS

		return true

	} else if value == "dns-over-tls" || value == "dot" {

		*protocol = ProtocolDNSoverTLS

		return true

	} else if value == "icmp" {

		*protocol = ProtocolICMP

		return true

	} else if value == "ssh" {

		*protocol = ProtocolSSH

		return true

	} else if value == "socks" || value == "socks4" || value == "socks5" {

		*protocol = ProtocolSOCKS

		return true

	} else if value == "whois" {

		*protocol = ProtocolWHOIS

		return true

	} else if value == "tcp" || value == "tcp4" || value == "tcp6" {

		*protocol = ProtocolTCP

		return true

	} else if value == "udp" || value == "udp4" || value == "udp6" {

		*protocol = ProtocolUDP

		return true

	}

	return false

}

func (protocol Protocol) Bytes() []byte {
	return []byte(protocol)
}

func (protocol Protocol) Equal(other Protocol) bool {
	return protocol == other
}

func (protocol Protocol) Hash() string {
	return toHash(protocol.String())
}

func (protocol Protocol) IsAny() bool {
	return protocol == "" || protocol == Protocol(Any)
}

func (protocol Protocol) IsValid() bool {

	switch protocol {
	case "", Protocol(Any), ProtocolHTTPS, ProtocolHTTP, ProtocolDNS, ProtocolDNSoverTLS, ProtocolICMP, ProtocolSSH, ProtocolSOCKS, ProtocolWHOIS, ProtocolTCP, ProtocolUDP:
		return true
	}

	return false

}

func (protocol Protocol) String() string {
	return string(protocol)
}

func (protocol Protocol) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(protocol))), nil
}

func (protocol *Protocol) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseProtocol(unquoted)

	if tmp != nil {
		*protocol = *tmp
	}

	return nil

}
