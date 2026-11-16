package types

import "math/big"
import "net/netip"
import "strings"

func addrToBig(value netip.Addr, bits int) *big.Int {

	var bytes []byte

	if bits == 32 {
		tmp := value.As4()
		bytes = tmp[:]
	} else {
		tmp := value.As16()
		bytes = tmp[:]
	}

	return new(big.Int).SetBytes(bytes)

}

func bigToAddr(value *big.Int, bits int) netip.Addr {

	var bytes []byte

	if bits == 32 {
		bytes = make([]byte, 4)
	} else {
		bytes = make([]byte, 16)
	}

	tmp := value.Bytes()

	if len(tmp) > len(bytes) {
		tmp = tmp[len(tmp)-len(bytes):]
	}

	copy(bytes[len(bytes)-len(tmp):], tmp)

	if bits == 32 {

		var tmp4 [4]byte
		copy(tmp4[:], bytes)

		return netip.AddrFrom4(tmp4)

	} else {

		var tmp16 [16]byte
		copy(tmp16[:], bytes)

		return netip.AddrFrom16(tmp16)

	}

}

func trailingZeros(value *big.Int, bits int) int {

	if value.Sign() == 0 {
		return bits
	}

	result := 0

	for result < bits && value.Bit(result) == 0 {
		result++
	}

	return result

}

func parseAddress(value string) (netip.Addr, bool) {

	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")

	tmp, err := netip.ParseAddr(value)

	if err == nil {
		return tmp.Unmap(), true
	}

	return netip.Addr{}, false

}

// ParseNetworks parses a CIDR prefix, a single IP address, or a
// "start - end" address range into one or more canonical prefixes.
func ParseNetworks(value string) []netip.Prefix {

	result := make([]netip.Prefix, 0)

	value = strings.TrimSpace(value)

	if value == "" {
		return result
	}

	if strings.Contains(value, " - ") {

		tmp := strings.Split(value, " - ")

		if len(tmp) == 2 {

			start, ok1 := parseAddress(tmp[0])
			end, ok2 := parseAddress(tmp[1])

			if ok1 == true && ok2 == true && start.Is4() == end.Is4() {

				bits := 128

				if start.Is4() {
					bits = 32
				}

				start_num := addrToBig(start, bits)
				end_num := addrToBig(end, bits)

				if start_num.Cmp(end_num) > 0 {
					start_num, end_num = end_num, start_num
				}

				one := big.NewInt(1)

				for start_num.Cmp(end_num) <= 0 {

					max_by_start := bits - trailingZeros(start_num, bits)

					span := new(big.Int).Sub(end_num, start_num)
					span.Add(span, one)

					size_bits := span.BitLen() - 1

					max_by_size := bits - size_bits
					prefix_length := max_by_start

					if max_by_size > prefix_length {
						prefix_length = max_by_size
					}

					prefix := netip.PrefixFrom(bigToAddr(start_num, bits), prefix_length)

					if prefix.IsValid() == true {
						result = append(result, prefix.Masked())
					}

					step := new(big.Int).Lsh(one, uint(bits-prefix_length))
					start_num.Add(start_num, step)

				}

			}

		}

	} else if strings.Contains(value, "/") {

		tmp, err := netip.ParsePrefix(value)

		if err == nil {
			result = append(result, tmp.Masked())
		}

	} else {

		tmp, ok := parseAddress(value)

		if ok == true {

			bits := 128

			if tmp.Is4() {
				bits = 32
			}

			result = append(result, netip.PrefixFrom(tmp, bits).Masked())

		}

	}

	return result

}

// NormalizeNetwork returns the first canonical prefix for a value.
func NormalizeNetwork(value string) string {

	tmp := ParseNetworks(value)

	if len(tmp) > 0 {
		return tmp[0].Masked().String()
	}

	return ""

}

// NormalizeNetworkList returns all canonical prefixes for a value. This is
// relevant for RIPE-style "start - end" ranges which are not prefix-based.
func NormalizeNetworkList(value string) []string {

	result := make([]string, 0)
	tmp := ParseNetworks(value)

	for p := 0; p < len(tmp); p++ {
		result = append(result, tmp[p].Masked().String())
	}

	return result

}

// NetworkContains checks whether an IP address is contained in a CIDR prefix
// or address range.
func NetworkContains(network string, value string) bool {

	tmp, ok := parseAddress(value)

	if ok == false {
		return false
	}

	prefixes := ParseNetworks(network)

	for p := 0; p < len(prefixes); p++ {

		if prefixes[p].Contains(tmp) == true {
			return true
		}

	}

	return false

}

// NetworkContainsPrefix checks whether inner is fully covered by outer.
func NetworkContainsPrefix(outer string, inner string) bool {

	outer_prefixes := ParseNetworks(outer)
	inner_prefixes := ParseNetworks(inner)

	for o := 0; o < len(outer_prefixes); o++ {

		for i := 0; i < len(inner_prefixes); i++ {

			if outer_prefixes[o].Addr().Is4() != inner_prefixes[i].Addr().Is4() {
				continue
			}

			if outer_prefixes[o].Bits() <= inner_prefixes[i].Bits() {

				if outer_prefixes[o].Contains(inner_prefixes[i].Addr()) == true {
					return true
				}

			}

		}

	}

	return false

}

// NetworkAddressCount returns the number of addresses in a prefix, capped at
// the maximum uint64 value for very large IPv6 prefixes.
func NetworkAddressCount(network string) uint64 {

	prefixes := ParseNetworks(network)

	if len(prefixes) == 0 {
		return 0
	}

	prefix := prefixes[0]
	bits := 128

	if prefix.Addr().Is4() {
		bits = 32
	}

	host_bits := bits - prefix.Bits()

	if host_bits >= 64 {
		return ^uint64(0)
	}

	return uint64(1) << uint(host_bits)

}

// ObservedPrefix returns a conservative prefix for an observed IP address.
func ObservedPrefix(value string) string {

	tmp, ok := parseAddress(value)

	if ok == false {
		return ""
	}

	if tmp.Is4() {
		return netip.PrefixFrom(tmp, 24).Masked().String()
	}

	return netip.PrefixFrom(tmp, 48).Masked().String()

}

// IsPrivateIP checks whether an IP address is not publicly routable.
func IsPrivateIP(value string) bool {

	tmp, ok := parseAddress(value)

	if ok == false {
		return false
	}

	if tmp.IsPrivate() == true || tmp.IsLoopback() == true || tmp.IsLinkLocalUnicast() == true || tmp.IsMulticast() == true || tmp.IsUnspecified() == true {
		return true
	}

	return false

}

// CanonicalIP returns a normalized IP address string, or the input if it
// cannot be parsed.
func CanonicalIP(value string) string {

	tmp, ok := parseAddress(value)

	if ok == true {
		return tmp.String()
	}

	return strings.TrimSpace(value)

}
