package dns

import "encoding/binary"
import "crypto/tls"
import "errors"
import "net"
import "strconv"

func resolveTLS(domain string, ip string, port uint16, query Packet) (Packet, error) {

	err := error(nil)
	response := NewPacket()
	server := net.ParseIP(ip)

	if server != nil {

		connection, err1 := tls.Dial("tcp", server.String()+":" + strconv.Itoa(int(port)), &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         domain,
		})

		if err1 == nil {

			query_buffer := query.Bytes()

			query_size := []byte{
				byte(len(query_buffer) >> 8),
				byte(len(query_buffer) & 0xff),
			}

			buffer := make([]byte, 0)
			buffer = append(buffer, query_size...)
			buffer = append(buffer, query_buffer...)

			_, err2 := connection.Write(buffer)

			if err2 == nil {

				response_buffer := make([]byte, 65535)
				response_size, err3 := connection.Read(response_buffer)

				if err3 == nil && response_size > 0 && response_size <= len(response_buffer) {

					response_buffer = response_buffer[0:response_size]
					packet_size := binary.BigEndian.Uint16(response_buffer[0:2])

					if packet_size > 0 && int(packet_size-2) <= len(response_buffer) {
						response = Parse(response_buffer[2 : packet_size+2])
					}

					connection.Close()

				} else {
					err = errors.New("Cannot parse buffer with length \"" + strconv.Itoa(response_size) + "\"")
					defer connection.Close()
				}

			} else {
				defer connection.Close()
				err = err2
			}

		} else {
			err = err1
		}

	} else {
		err = errors.New("Cannot parse IP \"" + ip + "\"")
	}

	return response, err

}
