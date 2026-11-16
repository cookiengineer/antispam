package whois

func parseBool(value string, thetype string) (bool, error) {

	val := false
	err := error(nil)

	if thetype == "bool" {

		if value == "true" {
			val = true
		}

	}

	return val, err

}
