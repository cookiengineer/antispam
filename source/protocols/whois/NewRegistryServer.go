package whois

import "antispam/types"
import "antispam/protocols/whois/iana"
import recon_types "antispam/types"
import "strings"

func NewRegistryServer(question string) *recon_types.Server {

	var result *recon_types.Server = nil

	if types.IsIPv4(question) {

		preflight_server := NewServer("whois.iana.org")
		preflight_schema := iana.IPv4{}
		preflight_buffer := Resolve(question, &preflight_server)

		if len(preflight_buffer) > 0 {

			preflight_err := Unmarshal(preflight_buffer, &preflight_schema)

			if preflight_err == nil {

				if preflight_schema.Whois != "" {
					tmp := NewServer(preflight_schema.Whois)
					result = &tmp
				}

			}

		}

	} else if types.IsIPv6(question) {

		preflight_server := NewServer("whois.iana.org")
		preflight_schema := iana.IPv6{}

		if strings.HasPrefix(question, "[") && strings.HasSuffix(question, "]") {
			question = question[1 : len(question)-1]
		}

		preflight_buffer := Resolve(question, &preflight_server)

		if len(preflight_buffer) > 0 {

			preflight_err := Unmarshal(preflight_buffer, &preflight_schema)

			if preflight_err == nil {

				if preflight_schema.Whois != "" {
					tmp := NewServer(preflight_schema.Whois)
					result = &tmp
				}

			}

		}

	} else if types.IsASN(question) {

		preflight_server := NewServer("whois.iana.org")
		preflight_schema := iana.ASN{}
		preflight_buffer := Resolve(question, &preflight_server)

		if len(preflight_buffer) > 0 {

			preflight_err := Unmarshal(preflight_buffer, &preflight_schema)

			if preflight_err == nil {

				if preflight_schema.Whois != "" {
					tmp := NewServer(preflight_schema.Whois)
					result = &tmp
				}

			}

		}

	}

	return result

}
