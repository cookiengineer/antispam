package types

import utils_net "antispam/utils/net"
import "math/rand"
import "strings"

type Server struct {
	Domain    string   `json:"domain"`
	Addresses []string `json:"addresses"`
	Port      uint16   `json:"port"`
	Protocol  Protocol `json:"protocol"`
	Schema    string   `json:"schema"`
}

func NewServer() Server {

	var server Server

	server.Domain = ""
	server.Port = 0
	server.Protocol = Any
	server.Schema = "DEFAULT"
	server.Addresses = make([]string, 0)

	return server

}

func (server *Server) IsIdentical(value Server) bool {

	result := false

	if server.Domain == value.Domain && server.Port == value.Port && server.Protocol == value.Protocol {
		result = true
	}

	return result

}

func (server *Server) IsValid() bool {

	result := false

	if server.Domain != "" && server.Port != 0 && server.Protocol != Any {

		if len(server.Addresses) > 0 {
			result = true
		}

	}

	return result

}

func (server *Server) AddAddress(value string) {

	valid := false

	if IsIPv6(value) {

		ipv6 := ParseIPv6(value)

		if ipv6 != nil {
			value = ipv6.String()
			valid = true
		}

	} else if IsIPv4(value) {

		ipv4 := ParseIPv4(value)

		if ipv4 != nil {
			value = ipv4.String()
			valid = true
		}

	}

	if valid == true {

		found := false

		for a := 0; a < len(server.Addresses); a++ {

			if server.Addresses[a] == value {
				found = true
				break
			}

		}

		if found == false {
			server.Addresses = append(server.Addresses, value)
		}

	}

}

func (server *Server) RemoveAddress(value string) {

	index := -1

	for a := 0; a < len(server.Addresses); a++ {

		if server.Addresses[a] == value {
			index = a
			break
		}

	}

	if index != -1 {
		server.Addresses = append(server.Addresses[:index], server.Addresses[index+1:]...)
	}

}

func (server *Server) SetAddresses(values []string) {

	filtered := make([]string, 0)

	for _, value := range values {

		if IsIPv6(value) {

			ipv6 := ParseIPv6(value)

			if ipv6 != nil {
				filtered = append(filtered, ipv6.String())
			}

		} else if IsIPv4(value) {

			ipv4 := ParseIPv4(value)

			if ipv4 != nil {
				filtered = append(filtered, ipv4.String())
			}

		}

	}

	server.Addresses = filtered

}

func (server *Server) SetDomain(value string) {

	if IsDomain(value) {

		domain := ParseDomain(value)

		if domain != nil {
			server.Domain = domain.String()
		}

	}

}

func (server *Server) SetPort(value uint16) {

	if value > 0 && value < 65535 {
		server.Port = value
	}

}

func (server *Server) SetProtocol(value Protocol) {
	server.Protocol = value
}

func (server *Server) SetSchema(value string) {
	server.Schema = strings.TrimSpace(value)
}

func (server *Server) RandomizeAddress() string {

	result := ""

	if utils_net.SupportsIPv4() {

		filtered := make([]string, 0)

		for a := 0; a < len(server.Addresses); a++ {

			address := server.Addresses[a]

			if IsIPv4(address) {
				filtered = append(filtered, address)
			}

		}

		if len(filtered) > 0 {
			result = filtered[rand.Intn(len(filtered))]
		}

	} else if utils_net.SupportsIPv6() {

		filtered := make([]string, 0)

		for a := 0; a < len(server.Addresses); a++ {

			address := server.Addresses[a]

			if IsIPv6(address) {
				filtered = append(filtered, address)
			}

		}

		if len(filtered) > 0 {
			result = filtered[rand.Intn(len(filtered))]
		}

	}

	return result

}
