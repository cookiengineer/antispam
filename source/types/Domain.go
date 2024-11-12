package types

import "strconv"
import "strings"

type Domain string

func isASCIILabel(value string) bool {

	var result bool = true

	for v := 0; v < len(value); v++ {

		character := string(value[v])

		if character >= "0" && character <= "9" {
			continue
		} else if character >= "A" && character <= "Z" {
			continue
		} else if character >= "a" && character <= "z" {
			continue
		} else if character == "-" || character == "_" {
			continue
		} else {
			result = false
			break
		}

	}

	return result

}

func IsDomain(value string) bool {

	if strings.Contains(value, ".") {

		labels := strings.Split(value, ".")
		valid := true

		if len(labels) >= 2 {

			for l := 0; l < len(labels)-1; l++ {

				label := labels[l]

				if isASCIILabel(label) && len(label) >= 1 {
					// Do Nothing
				} else {
					valid = false
					break
				}

			}

			tld := labels[len(labels)-1]

			if isASCIILabel(tld) && len(tld) >= 2 {
				// Do Nothing
			} else {
				valid = false
			}

		} else if len(labels) == 1 {

			hostname := labels[0]

			if isASCIILabel(hostname) && len(hostname) >= 3 {
				valid = true
			}

		}

		return valid

	}

	return false

}

func IsDomainAndPort(value string) bool {

	if strings.Contains(value, ":") {

		tmp := strings.Split(value, ":")

		if len(tmp) == 2 {

			_, err := strconv.ParseUint(tmp[1], 10, 16)

			if IsDomain(tmp[0]) && err == nil {
				return true
			}

		}

	}

	return false

}

func ParseDomain(value string) *Domain {

	var result *Domain = nil

	if IsDomain(value) {
		domain := Domain(value)
		result = &domain
	}

	return result

}

func ParseDomainAndPort(value string) (*Domain, uint16) {

	var result_domain *Domain = nil
	var result_port uint16 = 0

	if strings.Contains(value, ":") {

		tmp1 := strings.Split(value, ":")

		if len(tmp1) == 2 {

			tmp2 := ParseDomain(tmp1[0])
			num, err := strconv.ParseUint(tmp1[1], 10, 16)

			if tmp2 != nil && err == nil {
				result_domain = tmp2
				result_port = uint16(num)
			}

		}

	}

	return result_domain, result_port

}

func (domain Domain) Bytes() []byte {
	return []byte(domain)
}

func (domain Domain) Scope() string {
	return "public"
}

func (domain Domain) String() string {
	return string(domain)
}
