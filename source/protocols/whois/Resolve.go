package whois

import "antispam/types"
import "net"
import "strings"
import "time"

func Resolve(question string, server *types.Server) []byte {

	response_buffer := make([]byte, 0)

	if question != "" && server != nil {

		address := strings.Trim(server.RandomizeAddress(), "[]")
		address = net.JoinHostPort(address, "43")

		connection, err1 := net.DialTimeout("tcp", address, 15*time.Second)

		if err1 == nil {

			defer connection.Close()

			connection.SetDeadline(time.Now().Add(30 * time.Second))

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
				} else {
					break
				}

			}

		}

	}

	return response_buffer

}
