package whois

import "antispam/types"
import "antispam/protocols/dns"

func NewServer(domain string) types.Server {

	result := types.Server{}

	for s := 0; s < len(Servers); s++ {

		if Servers[s].Domain == domain {
			result = Servers[s]
			break
		}

	}

	if result.Domain != domain {

		dns_packet, dns_err := dns.Resolve(domain)

		if dns_err == nil && dns_packet.Type == "response" {

			result.SetDomain(domain)
			result.SetProtocol("whois")
			result.SetSchema("DEFAULT")

			for a := 0; a < len(dns_packet.Answers); a++ {

				answer := dns_packet.Answers[a]

				if answer.Type == dns.TypeA {
					result.AddAddress(answer.ToIPv4())
				} else if answer.Type == dns.TypeAAAA {
					result.AddAddress(answer.ToIPv6())
				}

			}

		}

	}

	return result

}
