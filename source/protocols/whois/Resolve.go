package whois

import "antispam/types"
import "io"
import "net"
import "strings"
import "time"

func Resolve(question string, server *types.Server) []byte {

	response_buffer := make([]byte, 0)

	if question != "" && server != nil {

		remote_address := net.TCPAddr{
			IP:   net.ParseIP(server.RandomizeAddress()),
			Port: 43,
		}

		connection, err1 := net.DialTCP("tcp", nil, &remote_address)

		if err1 == nil {

			request_buffer := make([]byte, 0)

			if server.Domain == "whois.arin.net" {

				if strings.HasPrefix(question, "AS") {

					request_buffer = append(request_buffer, []byte("a ")...)
					request_buffer = append(request_buffer, []byte(question[2:])...)
					request_buffer = append(request_buffer, []byte("\r\n")...)

				} else {

					request_buffer = append(request_buffer, []byte(question)...)
					request_buffer = append(request_buffer, []byte("\r\n")...)

				}

			} else if server.Domain == "whois.ripe.net" {

				request_buffer = append(request_buffer, []byte("-V Md5.5.18 ")...)
				request_buffer = append(request_buffer, []byte(question)...)
				request_buffer = append(request_buffer, []byte("\r\n")...)

			} else {

				request_buffer = append(request_buffer, []byte(question)...)
				request_buffer = append(request_buffer, []byte("\r\n")...)

			}

			connection.Write(request_buffer)

			response_buffer = make([]byte, 0)

			for {

				chunk := make([]byte, 65535)
				size, err2 := connection.Read(chunk)

				if err2 == nil && size > 0 {
					response_buffer = append(response_buffer, chunk[0:size]...)
				} else if err2 == io.EOF {
					break
				} else if size == 0 {
					time.Sleep(1 * time.Second)
				}

			}

			defer connection.Close()

		}

	}

	return response_buffer

}
