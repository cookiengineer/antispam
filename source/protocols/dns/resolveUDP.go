package dns

import "errors"
import "net"
import "strconv"
import "time"

func resolveUDP(ip string, port uint16, query Packet) (Packet, error) {

	err := error(nil)
	response := NewPacket()
	server := net.ParseIP(ip)

	if server != nil {

		udp_address := net.UDPAddr{
			IP:   server,
			Port: int(port),
		}

		connection, err1 := net.DialUDP("udp", nil, &udp_address)

		if err1 == nil {

			connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			connection.SetReadDeadline(time.Now().Add(10 * time.Second))

			_, err2 := connection.Write(query.Bytes())

			if err2 == nil {

				response_buffer := make([]byte, 1232)
				response_size, err3 := connection.Read(response_buffer)

				if err3 == nil && response_size > 0 && response_size < len(response_buffer) {
					response = Parse(response_buffer[0:response_size])
				} else {
					err = errors.New("Cannot parse buffer with length \"" + strconv.Itoa(response_size) + "\"")
				}

			} else {
				err = err2
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

