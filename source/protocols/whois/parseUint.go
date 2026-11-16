package whois

import "strconv"

func parseUint(value string, thetype string) (uint64, error) {

	num := uint64(0)
	err := error(nil)

	if thetype == "uint" {
		num, err = strconv.ParseUint(value, 10, 64)
	} else if thetype == "uint8" {
		num, err = strconv.ParseUint(value, 10, 8)
	} else if thetype == "uint16" {
		num, err = strconv.ParseUint(value, 10, 16)
	} else if thetype == "uint32" {
		num, err = strconv.ParseUint(value, 10, 32)
	} else if thetype == "uint64" {
		num, err = strconv.ParseUint(value, 10, 64)
	}

	return num, err

}
