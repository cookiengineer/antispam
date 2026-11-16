package dns

func parseLabels(buffer []byte, pointer *Pointer) ([]string, uint) {

	if pointer == nil {
		pointer = &Pointer{Offset: 0}
	}

	labels := make([]string, 0)
	bytes := uint(0)
	offset := pointer.Offset

	for offset < uint(len(buffer)) && bytes <= 255 {

		length := buffer[offset]

		if length == 0 {

			offset++
			bytes++
			break

		} else if length > 0 && length < 64 {

			label := string(buffer[offset+1 : offset+1+uint(length)])
			labels = append(labels, label)

			bytes += 1 + uint(length)
			offset += 1 + uint(length)

		} else if length > 64 {

			pointer_start := uint(((uint(length) - 0b11000000) << 8) + uint(buffer[offset+1]))
			pointer := Pointer{Offset: pointer_start}

			if pointer.Offset >= 12 && pointer.Offset < uint(len(buffer)) {

				pointer_labels, _ := parseLabels(buffer, &pointer)

				if len(pointer_labels) > 0 {

					labels = append(labels, pointer_labels...)
					bytes += 2
					offset += 2

					break

				}

			}

		}

	}

	return labels, bytes

}
