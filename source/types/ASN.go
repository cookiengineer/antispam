package types

import "strconv"
import "strings"

type ASN int

func IsASN(value string) bool {

	var result bool = false

	if strings.HasPrefix(value, "AS") && len(value) >= 3 {

		result = true

		for v := 2; v < len(value); v++ {

			chr := string(value[v])

			if chr >= "0" && chr <= "9" {
				continue
			} else {
				result = false
				break
			}

		}

	}

	return result

}

func NewASN() ASN {
	return ASN(0)
}

func ParseASN(value string) *ASN {

	var result *ASN = nil
	var asn ASN

	if asn.Parse(value) == true {
		result = &asn
	}

	return result

}

func ToASN(value string) ASN {

	var asn ASN

	asn.Parse(value)

	return asn

}

func (asn *ASN) Parse(value string) bool {

	if strings.HasPrefix(value, "AS") {

		num, err := strconv.ParseInt(value[2:], 10, 64)

		if err == nil {

			*asn = ASN(num)

			return true

		}

	} else {

		num, err := strconv.ParseInt(value, 10, 64)

		if err == nil {

			*asn = ASN(num)

			return true

		}

	}

	return false

}

func (asn ASN) Equal(other ASN) bool {
	return asn == other
}

func (asn ASN) Hash() string {
	return toHash(asn.String())
}

func (asn ASN) IsValid() bool {
	return asn >= 0
}

func (asn ASN) String() string {
	return "AS" + strconv.Itoa(int(asn))
}

func (asn ASN) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("AS" + strconv.Itoa(int(asn)))), nil
}

func (asn *ASN) UnmarshalJSON(data []byte) error {

	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	tmp := ParseASN(unquoted)

	if tmp != nil {
		*asn = *tmp
	}

	return nil

}
