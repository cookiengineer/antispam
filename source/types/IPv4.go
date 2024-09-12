package types

import "math"
import "strconv"
import "strings"

type IPv4 [4]byte

func IsIPv4(value string) bool {

	if strings.Contains(value, ".") {

		tmp := strings.Split(value, ".")

		if len(tmp) == 4 {

			valid := true

			for t := 0; t < len(tmp); t++ {

				_, err := strconv.ParseUint(tmp[t], 10, 8)

				if err != nil {
					valid = false
					break
				}

			}

			return valid

		}

	}

	return false

}

func IsIPv4AndPort(value string) bool {

	if strings.Contains(value, ".") && strings.Contains(value, ":") {

		tmp := strings.Split(value, ":")

		if len(tmp) == 2 {

			_, err := strconv.ParseUint(tmp[1], 10, 16)

			if IsIPv4(tmp[0]) && err == nil {
				return true
			}

		}

	}

	return false

}

func ParseIPv4(value string) *IPv4 {

	var result *IPv4 = nil

	if strings.Contains(value, ".") {

		tmp1 := strings.Split(value, ".")
		tmp2 := make([]uint8, 0)

		for t := 0; t < len(tmp1); t++ {

			num, err := strconv.ParseUint(tmp1[t], 10, 8)

			if err == nil {

				if num >= 0 && num <= 255 {
					tmp2 = append(tmp2, uint8(num))
				}

			}

		}

		if len(tmp2) == 4 {
			ipv4 := IPv4(tmp2)
			result = &ipv4
		}

	}

	return result

}

func ParseIPv4AndPort(value string) (*IPv4, uint16) {

	var result_ipv4 *IPv4 = nil
	var result_port uint16 = 0

	if strings.Contains(value, ".") && strings.Contains(value, ":") {

		tmp1 := strings.Split(value, ":")

		if len(tmp1) == 2 {

			tmp2 := ParseIPv4(tmp1[0])
			num, err := strconv.ParseUint(tmp1[1], 10, 16)

			if tmp2 != nil && err == nil {
				result_ipv4 = tmp2
				result_port = uint16(num)
			}

		}

	}

	return result_ipv4, result_port

}

func (ipv4 IPv4) Bytes(prefix uint8) []byte {

	bytes := make([]byte, 4)

	for i := 0; i < len(ipv4); i++ {
		bytes[i] = byte(ipv4[i])
	}

	if prefix%8 == 0 {

		bytes_length := int(prefix / 8)
		bytes_mask := make([]byte, bytes_length)
		copy(bytes_mask, bytes)

		for len(bytes_mask) < len(bytes) {
			bytes_mask = append(bytes_mask, uint8(0))
		}

		return bytes_mask

	} else {

		bytes_length := int(math.Floor(float64(prefix / 8)))
		bytes_mask := make([]byte, bytes_length)
		copy(bytes_mask, bytes)

		var last_byte byte = bytes[bytes_length]
		last_byte_shift := int(int(prefix) - bytes_length*8)
		bytes_mask = append(bytes_mask, last_byte>>(8-last_byte_shift))

		for len(bytes_mask) < len(bytes) {
			bytes_mask = append(bytes_mask, uint8(0))
		}

		return bytes_mask

	}

}

func (ipv4 IPv4) Scope() string {

	var result string = "public"

	private_ipv4s := []string{

		// RFC1122
		"0.0.0.0",
		"127.",

		// RFC1918
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
		"192.168.",

		// RFC2544
		"198.0.", "198.2.", "198.4.", "198.6.",
		"198.8.", "198.10.", "198.12.", "198.14.",
		"198.16.", "198.18.", "198.20.", "198.22.",
		"198.24.", "198.26.", "198.28.", "198.30.",
		"198.32.", "198.34.", "198.36.", "198.38.",
		"198.40.", "198.42.", "198.44.", "198.46.",
		"198.48.", "198.50.", "198.52.", "198.54.",
		"198.56.", "198.58.", "198.60.", "198.62.",
		"198.64.", "198.66.", "198.68.", "198.70.",
		"198.72.", "198.74.", "198.76.", "198.78.",
		"198.80.", "198.82.", "198.84.", "198.86.",
		"198.88.", "198.90.", "198.92.", "198.94.",
		"198.96.", "198.98.", "198.100.", "198.102.",
		"198.104.", "198.106.", "198.108.", "198.110.",
		"198.112.", "198.114.", "198.116.", "198.118.",
		"198.120.", "198.122.", "198.124.", "198.126.",
		"198.128.", "198.130.", "198.132.", "198.134.",
		"198.136.", "198.138.", "198.140.", "198.142.",
		"198.144.", "198.146.", "198.148.", "198.150.",
		"198.152.", "198.154.", "198.156.", "198.158.",
		"198.160.", "198.162.", "198.164.", "198.166.",
		"198.168.", "198.170.", "198.172.", "198.174.",
		"198.176.", "198.178.", "198.180.", "198.182.",
		"198.184.", "198.186.", "198.188.", "198.190.",
		"198.192.", "198.194.", "198.196.", "198.198.",
		"198.200.", "198.202.", "198.204.", "198.206.",
		"198.208.", "198.210.", "198.212.", "198.214.",
		"198.216.", "198.218.", "198.220.", "198.222.",
		"198.224.", "198.226.", "198.228.", "198.230.",
		"198.232.", "198.234.", "198.236.", "198.238.",
		"198.240.", "198.242.", "198.244.", "198.246.",
		"198.248.", "198.250.", "198.252.", "198.254.",

		// RFC3068
		"192.88.99.",

		// RFC3927
		"169.254.",

		// RFC6598
		"100.64.",
		"100.128.",
		"100.192.",

		// RFC6890
		"192.0.2.",
		"198.51.100.",
		"203.0.113.",
	}
	value := ipv4.String()

	for p := 0; p < len(private_ipv4s); p++ {

		if strings.HasPrefix(value, private_ipv4s[p]) {
			result = "private"
			break
		}

	}

	return result

}

func (ipv4 IPv4) String() string {

	tmp := make([]string, 4)

	for i := 0; i < len(ipv4); i++ {

		str := strconv.FormatUint(uint64(ipv4[i]), 10)

		if str != "" {
			tmp[i] = str
		} else {
			tmp[i] = "0"
		}

	}

	return strings.Join(tmp, ".")

}
