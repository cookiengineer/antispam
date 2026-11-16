package dns

func Parse(buffer []byte) Packet {

	packet := NewPacket()

	var pointer Pointer

	if len(buffer) >= 12 {

		packet.Identifier = uint16(buffer[0])<<8 + uint16(buffer[1])

		if buffer[2]&0b10000000 == 0 {
			packet.SetType("query")
		} else {
			packet.SetType("response")
		}

		packet.Codes.Operation = OperationCode(buffer[2] & 0b01111000)

		if buffer[2]&0b00000100 == 0 {
			packet.Flags.AuthorativeAnswer = false
		} else {
			packet.Flags.AuthorativeAnswer = true
		}

		if buffer[2]&0b00000010 == 0 {
			packet.Flags.Truncated = false
		} else {
			packet.Flags.Truncated = true
		}

		if buffer[2]&0b00000001 == 0 {
			packet.Flags.RecursionDesired = false
		} else {
			packet.Flags.RecursionDesired = true
		}

		if buffer[3]&0b10000000 == 0 {
			packet.Flags.RecursionAvailable = false
		} else {
			packet.Flags.RecursionAvailable = true
		}

		packet.Codes.Response = ResponseCode(buffer[3] & 0b00001111)

		qdcount := int(buffer[4])<<8 + int(buffer[5])
		ancount := int(buffer[6])<<8 + int(buffer[7])
		nscount := int(buffer[8])<<8 + int(buffer[9])
		arcount := int(buffer[10])<<8 + int(buffer[11])

		pointer.Offset = 12

		for q := 0; q < qdcount; q++ {

			question, err := parseQuestion(buffer, &pointer)

			if err == nil {
				packet.Questions = append(packet.Questions, *question)
			}

		}

		for a := 0; a < ancount; a++ {

			answer, err := parseRecord(buffer, &pointer)

			if err == nil {
				packet.Answers = append(packet.Answers, *answer)
			}

		}

		for n := 0; n < nscount; n++ {

			authority, err := parseRecord(buffer, &pointer)

			if err == nil {
				packet.Authorities = append(packet.Authorities, *authority)
			}

		}

		for a := 0; a < arcount; a++ {

			record, err := parseRecord(buffer, &pointer)

			if err == nil {
				packet.Additionals = append(packet.Additionals, *record)
			}

		}

	}

	return packet

}
