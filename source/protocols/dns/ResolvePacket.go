package dns

import "antispam/types"
import recon_types "antispam/types"
import "errors"
import "math/rand"
import "time"

func ResolvePacket(query Packet) (Packet, error) {

	response := Packet{}
	response_err := error(nil)

	if query.Server == nil {
		rand.Seed(time.Now().Unix())
		query.SetServer(Servers[rand.Intn(len(Servers))])
	}

	if query.Server != nil {

		tmp := Packet{}
		tmp_err := error(nil)

		if query.Server.Protocol == types.ProtocolDNS {

			has_special_questions := false

			for q := 0; q < len(query.Questions); q++ {

				question := query.Questions[q]

				if question.Type != TypeA && question.Type != TypeAAAA {
					has_special_questions = true
					break
				}

			}

			if has_special_questions == true {
				tmp, tmp_err = resolveUDP(query.Server.RandomizeAddress(), query.Server.Port, query)
			} else {
				tmp, tmp_err = resolveUDPBatch(query.Server.RandomizeAddress(), query.Server.Port, query)
			}

		} else if query.Server.Protocol == types.ProtocolDNSoverTLS {
			tmp, tmp_err = resolveTLS(query.Server.Domain, query.Server.RandomizeAddress(), query.Server.Port, query)
		} else {
			tmp_err = errors.New("Cannot resolve via protocol \"" + query.Server.Protocol.String() + "\"")
		}

		if tmp_err == nil && tmp.Type == "response" {

			if tmp.Flags.RecursionAvailable == true {

				response = tmp

			} else if len(tmp.Authorities) > 0 {

				nameserver_response := Packet{}
				nameserver_err := error(nil)

				nameserver_domain := ""
				nameserver_addresses := make([]string, 0)

				for a := 0; a < len(tmp.Authorities); a++ {

					domain := tmp.Authorities[a].ToDomain()

					if domain != "" {
						nameserver_domain = domain
						break
					}

				}

				if nameserver_domain != "" {

					nameserver_query := NewPacket()
					nameserver_query.SetType("query")
					nameserver_query.AddQuestion(NewQuestion(nameserver_domain, TypeA))
					nameserver_query.AddQuestion(NewQuestion(nameserver_domain, TypeAAAA))

					if query.Server.Protocol == types.ProtocolDNS {
						nameserver_response, nameserver_err = resolveUDP(query.Server.RandomizeAddress(), query.Server.Port, nameserver_query)
					} else if query.Server.Protocol == types.ProtocolDNSoverTLS {
						nameserver_response, nameserver_err = resolveTLS(query.Server.Domain, query.Server.RandomizeAddress(), query.Server.Port, nameserver_query)
					}

					if nameserver_err == nil && nameserver_response.Type == "response" {

						if len(nameserver_response.Answers) > 0 {

							for a := 0; a < len(nameserver_response.Answers); a++ {

								record := nameserver_response.Answers[a]

								if record.Type == TypeA {
									nameserver_addresses = append(nameserver_addresses, record.ToIPv4())
								} else if record.Type == TypeAAAA {
									nameserver_addresses = append(nameserver_addresses, record.ToIPv6())
								}

							}

						}

					}

				}

				if nameserver_domain != "" && len(nameserver_addresses) > 0 {

					nameserver := recon_types.Server{
						Domain:    nameserver_domain,
						Addresses: nameserver_addresses,
						Port:      53,
						Protocol:  types.ProtocolDNS,
						Schema:    "",
					}

					resolved_response, resolved_err := resolveUDP(nameserver.RandomizeAddress(), nameserver.Port, query)

					if resolved_err == nil && resolved_response.Type == "response" && resolved_response.Flags.RecursionAvailable == true {

						response = resolved_response
						response.SetIdentifier(query.Identifier)
						response.SetServer(nameserver)

					} else {
						// Do Nothing, probably a DNS Loop
						response_err = errors.New("ResolvePacket: Invalid Nameserver DNS response")
					}

				}

			} else {
				response_err = errors.New("ResolvePacket: Invalid DNS response")
			}

		} else {
			response_err = tmp_err
		}

	}

	return response, response_err

}
