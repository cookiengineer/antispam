package dns

var blocked_A []string = []string{
	"0.0.0.0",
	"1.1.1.3",
	"1.1.1.4",
	"1.1.1.5",
	"1.1.1.6",
	"1.1.1.7",
	"1.1.1.8",
	"1.1.1.9",
}

var blocked_AAAA []string = []string{
	"[0000:0000:0000:0000:0000:0000:0000:0000]",
	"[2606:4700:0000:0000:0000:0000:0000:1112]",
	"[2606:4700:0000:0000:0000:0000:0000:1113]",
	"[2606:4700:0000:0000:0000:0000:0000:1114]",
	"[2606:4700:0000:0000:0000:0000:0000:1115]",
	"[2606:4700:0000:0000:0000:0000:0000:1116]",
	"[2606:4700:0000:0000:0000:0000:0000:1117]",
	"[2606:4700:0000:0000:0000:0000:0000:1118]",
	"[2606:4700:0000:0000:0000:0000:0000:1119]",
}

func IsBlocked(response Packet) bool {

	result := false

	for a1 := 0; a1 < len(response.Answers); a1++ {

		answer := response.Answers[a1]

		if answer.Type == TypeA {

			for a2 := 0; a2 < len(blocked_A); a2++ {

				if answer.ToIPv4() == blocked_A[a2] {
					result = true
					break
				}

			}

		} else if answer.Type == TypeAAAA {

			for a2 := 0; a2 < len(blocked_AAAA); a2++ {

				if answer.ToIPv6() == blocked_AAAA[a2] {
					result = true
					break
				}

			}

		}

		if result == true {
			break
		}

	}

	return result

}
