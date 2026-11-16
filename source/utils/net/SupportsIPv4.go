package net

import "net"

func SupportsIPv4() bool {

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

						address_v4 := value.IP.To4()

						if address_v4 != nil {
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
