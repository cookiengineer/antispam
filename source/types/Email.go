package types

import utils_strings "antispam/utils/strings"
import "slices"
import "strings"
import "time"

func parseBoundary(content string) string {

	var result string

	lines := strings.Split(strings.TrimSpace(content), "\n")

	for l := 0; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])

		if strings.HasPrefix(line, "boundary=") && strings.HasSuffix(line, ";") {

			result = utils_strings.TrimQuotes(strings.TrimSpace(line[9:len(line)-1]))
			break

		} else if strings.HasPrefix(line, "boundary=") {

			result = utils_strings.TrimQuotes(strings.TrimSpace(line[9:]))
			break

		} else if strings.HasSuffix(line, ";") {
			// Do Nothing
		} else if strings.Contains(line, ";") {

			chunks := strings.Split(line, ";")

			for c := 0; c < len(chunks); c++ {

				chunk := strings.TrimSpace(chunks[c])

				if strings.HasPrefix(chunk, "boundary=") {
					result = utils_strings.TrimQuotes(strings.TrimSpace(chunk[9:]))
					break
				}

			}

		}

		if strings.HasPrefix(line, "--") {
			break
		}

	}

	return result

}

type Email struct {
	MessageID string    `json:"message_id"`
	Boundary  string    `json:"boundary"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Date      time.Time `json:"date"`
	Domains   []string  `json:"domains"`
	IPv4s     []string  `json:"ipv4s"`
	IPv6s     []string  `json:"ipv6s"`
}

func NewEmail() Email {

	var email Email

	email.Domains = make([]string, 0)
	email.IPv4s = make([]string, 0)
	email.IPv6s = make([]string, 0)

	return email

}

func IsEmail(buffer []byte) bool {

	lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")

	found_from := false
	found_to := false
	found_received := false

	for l := 0; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])

		if strings.HasPrefix(line, "From: ") {
			found_from = true
		} else if strings.HasPrefix(line, "To: ") {
			found_to = true
		} else if strings.HasPrefix(line, "Received: ") {
			found_received = true
		}

	}

	return found_from && found_to && found_received

}

func ParseEmail(buffer []byte) *Email {

	var result *Email = nil

	email := NewEmail()
	lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")


	header := ""
	headers := make(map[string][]string, 0)

	for l := 0; l < len(lines); l++ {

		line := lines[l]

		if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {

			header = strings.ToLower(strings.TrimSpace(line[0:strings.Index(line, ":")]))
			headers[header] = []string{line[strings.Index(line, ":")+1:]}

		} else if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {

			if header != "" {
				headers[header] = append(headers[header], line)
			}

		} else if strings.HasPrefix(line, "--") {

			// MIME boundary
			header = ""
			break

		}

	}

	if _, ok := headers["message-id"]; ok {

		for _, raw := range headers["message-id"] {

			tmp := strings.TrimSpace(raw)

			if strings.HasPrefix(tmp, "<") && strings.HasSuffix(tmp, ">") {
				email.MessageID = tmp[1:len(tmp)-1]
				break
			}

		}

	}

	if _, ok := headers["content-type"]; ok {

		for _, raw := range headers["content-type"] {

			line := strings.TrimSpace(raw)

			boundary := parseBoundary(line)

			if boundary != "" {
				email.Boundary = boundary
			}

		}

	}

	if _, ok := headers["date"]; ok {

		for _, raw := range headers["date"] {

			tmp := strings.TrimSpace(raw)

			time1, err1 := time.Parse(time.RFC1123, tmp)
			time2, err2 := time.Parse(time.RFC1123Z, tmp)

			if err1 == nil {
				email.Date = time1
				break
			} else if err2 == nil {
				email.Date = time2
				break
			}

		}

	}

	if _, ok := headers["from"]; ok {

		for _, raw := range headers["from"] {

			tmp := strings.TrimSpace(raw)

			if strings.Contains(tmp, "<") && strings.HasSuffix(tmp, ">") {
				tmp = tmp[strings.Index(tmp, "<")+1:len(tmp)-1]
			}

			if strings.Contains(tmp, "@") {

				address := strings.ToLower(strings.TrimSpace(tmp))
				domain := strings.ToLower(strings.TrimSpace(address[strings.Index(address, "@")+1:]))

				email.From = address

				if !slices.Contains(email.Domains, domain) {
					email.Domains = append(email.Domains, domain)
				}

				break

			}

		}

	}

	if _, ok := headers["received"]; ok {

		for _, raw := range headers["received"] {

			line := strings.TrimSpace(raw)

			if strings.Contains(line, "by ") {
				line = strings.TrimSpace(line[0:strings.Index(line, "by ")])
			}

			if strings.HasPrefix(line, "from ") {

				tmp1 := strings.TrimSpace(line[5:])

				if strings.Contains(tmp1, " ") {

					domain := strings.ToLower(strings.TrimSpace(tmp1[0:strings.Index(tmp1, " ")]))

					if !slices.Contains(email.Domains, domain) {
						email.Domains = append(email.Domains, domain)
					}

				}

				if strings.Contains(tmp1, "(") && strings.Contains(tmp1, ")") {

					chunk := utils_strings.Cut(tmp1, "(", ")")

					if strings.Contains(chunk, " ") {

						relay_domain := strings.ToLower(strings.TrimSpace(chunk[0:strings.Index(chunk, " ")]))

						if strings.Contains(relay_domain, ".") {

							if !slices.Contains(email.Domains, relay_domain) {
								email.Domains = append(email.Domains, relay_domain)
							}

						}

					}

					if strings.Contains(tmp1, "[") && strings.Contains(tmp1, "]") {

						chunk_ip := utils_strings.Cut(tmp1, "[", "]")

						if IsIPv6(chunk_ip) {

							ipv6 := ParseIPv6(chunk_ip)

							if ipv6 != nil {

								tmp := ipv6.String()

								if !slices.Contains(email.IPv6s, tmp) {
									email.IPv6s = append(email.IPv6s, tmp)
								}

							}

						} else if IsIPv4(chunk_ip) {

							ipv4 := ParseIPv4(chunk_ip)

							if ipv4 != nil {

								tmp := ipv4.String()

								if !slices.Contains(email.IPv4s, tmp) {
									email.IPv4s = append(email.IPv4s, tmp)
								}

							}

						}

					} else if IsIPv6(chunk) {

						ipv6 := ParseIPv6(chunk)

						if ipv6 != nil {

							tmp := ipv6.String()

							if !slices.Contains(email.IPv6s, tmp) {
								email.IPv6s = append(email.IPv6s, tmp)
							}

						}

					} else if IsIPv4(chunk) {

						ipv4 := ParseIPv4(chunk)

						if ipv4 != nil {

							tmp := ipv4.String()

							if !slices.Contains(email.IPv4s, tmp) {
								email.IPv4s = append(email.IPv4s, tmp)
							}

						}

					}


				}

			}

		}

	}

	if _, ok := headers["return-path"]; ok {

		for _, raw := range headers["return-path"] {

			address := strings.TrimSpace(raw)

			if strings.HasPrefix(address, "<") && strings.HasSuffix(address, ">") {
				address = address[1:len(address)-1]
			}

			if strings.Contains(address, "@") {

				domain := strings.ToLower(strings.TrimSpace(address[strings.Index(address, "@")+1:]))

				if domain != "" {
					email.Domains = append(email.Domains, domain)
				}

			}

		}

	}

	if _, ok := headers["subject"]; ok {

		for _, raw := range headers["subject"] {

			tmp := strings.TrimSpace(raw)

			if tmp != "" {

				email.Subject = utils_strings.ToASCII(tmp)
				break

			}

		}

	}

	if _, ok := headers["to"]; ok {

		for _, raw := range headers["to"] {

			tmp := strings.TrimSpace(raw)

			if strings.Contains(tmp, "<") && strings.HasSuffix(tmp, ">") {
				tmp = tmp[strings.Index(tmp, "<")+1:len(tmp)-1]
			}

			if strings.Contains(tmp, "@") {

				address := strings.ToLower(strings.TrimSpace(tmp))
				domain := strings.ToLower(strings.TrimSpace(address[strings.Index(address, "@")+1:]))

				email.To = address

				if !slices.Contains(email.Domains, domain) {
					email.Domains = append(email.Domains, domain)
				}

				break

			}

		}

	}

	if email.Boundary != "" {

		raw_buffer := strings.TrimSpace(string(buffer))
		message := strings.TrimSpace(raw_buffer[strings.Index(raw_buffer, "\n--" + email.Boundary) + len(email.Boundary)+3:])

		if strings.Contains(message, email.Boundary) {

			contents := strings.Split(message, "\n--" + email.Boundary + "\n")

			for _, content := range contents {

				if strings.Contains(content, "Content-Type: multipart/alternative") {

					nested_boundary := parseBoundary(content)

					if nested_boundary != "" {

						// Outlook format, which uses nested MIME boundaries
						nested_contents := strings.Split(message, "\n--" + nested_boundary + "\n")

						for _, nested_content := range nested_contents {

							if strings.Contains(nested_content, "Content-Type: text/plain") {
								email.Message = strings.TrimSpace(nested_content)
								break
							}

						}

					}

				} else if strings.Contains(content, "Content-Type: text/plain") {

					// Apple Mail / Thunderbird format
					email.Message = strings.TrimSpace(content)
					break

				}

			}

		}

	}

	if email.To != "" && len(email.Domains) > 0 {

		email_domain := email.To[strings.Index(email.To, "@")+1:]
		filtered := make([]string, 0)

		for d := 0; d < len(email.Domains); d++ {

			domain := email.Domains[d]

			if domain == email_domain || strings.HasSuffix(domain, "." + email_domain) {
				// Do Nothing
			} else {
				filtered = append(filtered, domain)
			}

		}

		email.Domains = filtered

	}

	if email.MessageID != "" {
		result = &email
	}

	return result

}
