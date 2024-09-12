package strings

import "strings"

func ToASCIIName(value string) string {

	tmp := ToASCII(value)

	var filtered string
	var last_was_dash bool = false

	for _, chunk := range tmp {

		if chunk == 45 || chunk == 58 || chunk == 59 || chunk == 95 {

			if last_was_dash == false {
				filtered += "-"
				last_was_dash = true
			}

		} else if chunk >= 33 && chunk <= 47 {

			// Do Nothing

		} else if chunk >= 48 && chunk <= 57 {

			// Digits
			filtered += string(chunk)
			last_was_dash = false

		} else if chunk >= 58 && chunk <= 64 {

			// Do Nothing

		} else if chunk >= 65 && chunk <= 90 {

			// Uppercase letters
			filtered += string(chunk)
			last_was_dash = false

		} else if chunk >= 91 && chunk <= 96 {

			// Do Nothing

		} else if chunk >= 97 && chunk <= 122 {

			// Lowercase letters
			filtered += string(chunk)
			last_was_dash = false

		} else if chunk >= 123 && chunk <= 126 {

			// Do Nothing

		}

	}

	if filtered != "-" {
		filtered = strings.TrimSpace(filtered)
	}

	return filtered

}
