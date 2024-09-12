package types

import "math"
import "strconv"
import "strings"

func formatIPv6(value string) string {

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = value[1 : len(value)-1]
	}

	if strings.Contains(value, "::") {

		tmp := strings.Split(value, "::")
		prefix := strings.Split(tmp[0], ":")
		suffix := strings.Split(tmp[1], ":")
		complete := []string{"0000", "0000", "0000", "0000", "0000", "0000", "0000", "0000"}

		for p := 0; p < len(prefix); p++ {

			chunk := prefix[p]
			c := p

			if len(chunk) == 1 {
				complete[c] = "000" + chunk
			} else if len(chunk) == 2 {
				complete[c] = "00" + chunk
			} else if len(chunk) == 3 {
				complete[c] = "0" + chunk
			} else if len(chunk) == 4 {
				complete[c] = chunk
			}

		}

		for s := 0; s < len(suffix); s++ {

			chunk := suffix[s]
			c := 8 - len(suffix) + s

			if len(chunk) == 1 {
				complete[c] = "000" + chunk
			} else if len(chunk) == 2 {
				complete[c] = "00" + chunk
			} else if len(chunk) == 3 {
				complete[c] = "0" + chunk
			} else if len(chunk) == 4 {
				complete[c] = chunk
			}

		}

		return "[" + strings.Join(complete, ":") + "]"

	} else if strings.Contains(value, ":") {

		chunks := strings.Split(value, ":")
		complete := []string{"0000", "0000", "0000", "0000", "0000", "0000", "0000", "0000"}

		for c := 0; c < len(chunks); c++ {

			chunk := chunks[c]

			if len(chunk) == 1 {
				complete[c] = "000" + chunk
			} else if len(chunk) == 2 {
				complete[c] = "00" + chunk
			} else if len(chunk) == 3 {
				complete[c] = "0" + chunk
			} else if len(chunk) == 4 {
				complete[c] = chunk
			}

		}

		return "[" + strings.Join(complete, ":") + "]"

	}

	return ""

}

type IPv6 [16]byte

func IsIPv6(raw string) bool {

	value := formatIPv6(raw)

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {

		// Ignore embedded IPv6::IPv4 syntax
		if !strings.Contains(value[1:len(value)-1], ".") {

			tmp := strings.Split(value[1:len(value)-1], ":")

			if len(tmp) == 8 {

				valid := true

				for t := 0; t < len(tmp); t++ {

					_, err := strconv.ParseUint(tmp[t], 16, 64)

					if err != nil {
						valid = false
						break
					}

				}

				return valid

			}

		}

	} else if strings.Contains(value, ":") {

		// Ignore embedded IPv6::IPv4 syntax
		if !strings.Contains(value[1:len(value)-1], ".") {

			tmp := strings.Split(value, ":")

			if len(tmp) == 8 {

				valid := true

				for t := 0; t < len(tmp); t++ {

					_, err := strconv.ParseUint(tmp[t], 16, 64)

					if err != nil {
						valid = false
						break
					}

				}

				return valid

			}

		}

	}

	return false

}

func IsIPv6AndPort(value string) bool {

	if strings.HasPrefix(value, "[") && strings.Contains(value, "]:") {

		tmp := strings.Split(value[1:], "]:")

		if len(tmp) == 2 {

			_, err := strconv.ParseUint(tmp[1], 10, 16)

			if IsIPv6("[" + value + "]") && err == nil {
				return true
			}

		}

	}

	return false

}

func ParseIPv6(raw string) *IPv6 {

	var result *IPv6 = nil

	value := formatIPv6(raw)

	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {

		tmp1 := strings.Join(strings.Split(value[1:len(value)-1], ":"), "")
		tmp2 := make([]uint8, 0)

		for t := 0; t < len(tmp1); t += 2 {

			hex := string(tmp1[t : t+2])
			num, err := strconv.ParseUint(hex, 16, 8)

			if err == nil {
				tmp2 = append(tmp2, uint8(num))
			} else {
				tmp2 = append(tmp2, uint8(0))
			}

		}

		if len(tmp2) == 16 {
			ipv6 := IPv6(tmp2)
			result = &ipv6
		}

	}

	return result

}

func ParseIPv6AndPort(value string) (*IPv6, uint16) {

	var result_ipv6 *IPv6 = nil
	var result_port uint16 = 0

	if strings.HasPrefix(value, "[") && strings.Contains(value, "]:") {

		tmp1 := strings.Split(value[1:], "]:")

		if len(tmp1) == 2 {

			tmp2 := ParseIPv6("[" + tmp1[0] + "]")
			num, err := strconv.ParseUint(tmp1[1], 10, 16)

			if tmp2 != nil && err == nil {
				result_ipv6 = tmp2
				result_port = uint16(num)
			}

		}

	}

	return result_ipv6, result_port

}

func (ipv6 IPv6) Bytes(prefix uint8) []byte {

	bytes := make([]byte, 16)

	for i := 0; i < len(ipv6); i++ {
		bytes[i] = byte(ipv6[i])
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

func (ipv6 IPv6) Scope() string {

	var result string = "public"

	private_ipv6s := []string{

		// RFC3513
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"0000:0000:0000:0000:0000:0000:0000:0001",
		"fe80:0000:0000:0000",
	}
	value := ipv6.String()

	for p := 0; p < len(private_ipv6s); p++ {

		if strings.HasPrefix(value, "[" + private_ipv6s[p]) {
			result = "private"
			break
		}

	}

	return result

}

func (ipv6 IPv6) String() string {

	tmp := make([]string, 8)

	for i := 0; i < len(ipv6); i += 2 {

		hex1 := strconv.FormatUint(uint64(ipv6[i]), 16)
		hex2 := strconv.FormatUint(uint64(ipv6[i+1]), 16)

		if len(hex1) == 1 {
			hex1 = "0" + hex1
		}

		if len(hex2) == 1 {
			hex2 = "0" + hex2
		}

		tmp[i / 2] = hex1 + hex2

	}

	return strings.Join(tmp, ":")

}
