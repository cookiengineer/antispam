package dns

func renderLabels(labels []string) []byte {

	buffer := make([]byte, 0)

	if len(labels) == 0 {

		buffer = append(buffer, 0)
		buffer = append(buffer, 0)

	} else {

		for l := 0; l < len(labels); l++ {

			label := []byte(labels[l])
			length := uint8(len(label))

			buffer = append(buffer, byte(length))
			buffer = append(buffer, label...)

		}

		buffer = append(buffer, 0)

	}

	return buffer

}
