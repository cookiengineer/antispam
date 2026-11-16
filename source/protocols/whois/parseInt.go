package whois

import "strconv"

func parseInt(value string, thetype string) (int64, error) {

	num := int64(0)
	err := error(nil)

	if thetype == "int" {
		num, err = strconv.ParseInt(value, 10, 64)
	} else if thetype == "int8" {
		num, err = strconv.ParseInt(value, 10, 8)
	} else if thetype == "int16" {
		num, err = strconv.ParseInt(value, 10, 16)
	} else if thetype == "int32" {
		num, err = strconv.ParseInt(value, 10, 32)
	} else if thetype == "int64" {
		num, err = strconv.ParseInt(value, 10, 64)
	}

	return num, err

}
