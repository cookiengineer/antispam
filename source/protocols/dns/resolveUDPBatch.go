package dns

import "errors"
import "net"
import "strconv"
import "time"

func cloneIntoSingleQuestions(query Packet) []Packet {

	clones := make([]Packet, 0)

	if len(query.Questions) > 1 {

		for q := 0; q < len(query.Questions); q++ {

			clone := NewPacket()
			clone.SetType("query")
			clone.AddQuestion(query.Questions[q])

			clones = append(clones, clone)

		}

	} else {
		clones = append(clones, query)
	}

	return clones

}

func mergeIntoSingleResponse(response *Packet, partial Packet) {

	if response.Type != partial.Type {

		response.Type = partial.Type
		response.Server = partial.Server
		response.Codes.Operation = partial.Codes.Operation
		response.Codes.Response = partial.Codes.Response
		response.Flags.AuthorativeAnswer = partial.Flags.AuthorativeAnswer
		response.Flags.Truncated = partial.Flags.Truncated
		response.Flags.RecursionAvailable = partial.Flags.RecursionAvailable
		response.Flags.RecursionDesired = partial.Flags.RecursionDesired

	}

	for a := 0; a < len(partial.Answers); a++ {
		response.AddAnswer(partial.Answers[a])
	}

	for a := 0; a < len(partial.Authorities); a++ {
		response.AddAuthority(partial.Authorities[a])
	}

	for a := 0; a < len(partial.Additionals); a++ {
		response.AddAdditional(partial.Additionals[a])
	}

}

func resolveUDPBatch(ip string, port uint16, query Packet) (Packet, error) {

	err := error(nil)
	response := NewPacket()
	server := net.ParseIP(ip)

	if server != nil {

		for q := 0; q < len(query.Questions); q++ {
			response.AddQuestion(query.Questions[q])
		}

		udp_address := net.UDPAddr{
			IP:   server,
			Port: int(port),
		}

		packets := cloneIntoSingleQuestions(query)

		connection, err1 := net.DialUDP("udp", nil, &udp_address)

		if err1 == nil {

			connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			connection.SetReadDeadline(time.Now().Add(10 * time.Second))

			for p := 0; p < len(packets); p++ {

				_, err2 := connection.Write(packets[p].Bytes())

				if err2 == nil {

					response_buffer := make([]byte, 1232)
					response_size, err3 := connection.Read(response_buffer)

					if err3 == nil && response_size > 0 && response_size < len(response_buffer) {

						partial_response := Parse(response_buffer[0:response_size])
						mergeIntoSingleResponse(&response, partial_response)

					} else {
						err = errors.New("Cannot parse buffer with length \"" + strconv.Itoa(response_size) + "\"")
						break
					}

				} else {
					err = err2
					break
				}

			}

			defer connection.Close()

		} else {
			err = err1
		}

	} else {
		err = errors.New("Cannot parse IP \"" + ip + "\"")
	}

	return response, err

}

