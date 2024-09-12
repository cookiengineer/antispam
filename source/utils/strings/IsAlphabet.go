package strings

func IsAlphabet(value string) bool {

	var result bool = true

	for v := 0; v < len(value); v++ {

		character := string(value[v])

		if character >= "0" && character <= "9" {
			continue
		} else if character >= "A" && character <= "Z" {
			continue
		} else if character >= "a" && character <= "z" {
			continue
		} else {
			result = false
			break
		}

	}

	return result

}
