package whois

import "strconv"

func parseFloat(value string, thetype string) (float64, error) {

	num := float64(0.0)
	err := error(nil)

	if thetype == "float32" {
		num, err = strconv.ParseFloat(value, 32)
	} else if thetype == "float64" {
		num, err = strconv.ParseFloat(value, 64)
	}

	return num, err

}
