package net

import "net"

func SupportsIPv6() bool {

	result := false
	ifaces, err1 := net.Interfaces()

	if err1 == nil {

		for _, iface := range ifaces {

			addresses, err2 := iface.Addrs()

			if err2 == nil {

				for _, raw := range addresses {

					switch value := raw.(type) {

					case *net.IPAddr:

						// Do Nothing

					case *net.IPNet:

						address_v6 := value.IP.To16()

						if address_v6 != nil {
							result = true
							break
						}

					}

				}

			}

			if result == true {
				break
			}

		}

	}

	return result

}
