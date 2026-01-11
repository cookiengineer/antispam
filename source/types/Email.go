package types

import utils_strings "antispam/utils/strings"
import "encoding/base64"
import "slices"
import "strings"
import "time"

func toASCIIMessage(raw string) string {

	lines := make([]string, 0)

	raw_lines := strings.Split(strings.TrimSpace(string(raw)), "\n")

	for l := 0; l < len(raw_lines); l++ {

		line := utils_strings.ToASCII(raw_lines[l])
		lines = append(lines, line)

	}

	return strings.Join(lines, "\n")

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
		} else if strings.HasPrefix(line, "To: ") || strings.HasPrefix(line, "Delivered-To: ") {
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
	headers := make(map[string]string, 0)

	for _, line := range lines {

		if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {

			header = strings.ToLower(strings.TrimSpace(line[0:strings.Index(line, ":")]))

			_, is_multiline := headers[header];

			if is_multiline == true {

				tmp1 := strings.TrimSpace(headers[header])
				tmp2 := strings.TrimSpace(line[strings.Index(line, ":")+1:])

				headers[header] = tmp1 + "\n" + tmp2

			} else {

				tmp := strings.TrimSpace(line[strings.Index(line, ":")+1:])

				if tmp != "" {
					headers[header] = tmp
				}

			}

		} else if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {

			if header != "" {

				tmp := strings.TrimSpace(line)

				if tmp != "" {
					headers[header] = headers[header] + " " + tmp
				}

			}

		} else if strings.HasPrefix(line, "--") {

			// MIME boundary
			header = ""
			break

		}

	}

	if raw, ok := headers["message-id"]; ok {

		tmp := strings.TrimSpace(raw)

		if strings.HasPrefix(tmp, "<") && strings.HasSuffix(tmp, ">") {
			email.MessageID = tmp[1:len(tmp)-1]
		}

	}

	if raw, ok := headers["content-type"]; ok {

		tmp := strings.TrimSpace(raw)
		boundary := parseBoundary(tmp)

		if boundary != "" {
			email.Boundary = boundary
		}

	}

	if raw, ok := headers["date"]; ok {

		tmp := strings.TrimSpace(raw)

		time1, err1 := time.Parse(time.RFC1123, tmp)
		time2, err2 := time.Parse(time.RFC1123Z, tmp)

		if err1 == nil {
			email.Date = time1
		} else if err2 == nil {
			email.Date = time2
		}

	}

	if raw, ok := headers["received"]; ok {

		tmp := strings.TrimSpace(raw)
		trace_path := parseReceived(tmp)

		for _, host := range trace_path {

			if IsIPv6(host) {

				ipv6 := ParseIPv6(host)

				if ipv6 != nil {

					tmp := ipv6.String()

					if !slices.Contains(email.IPv6s, tmp) {
						email.IPv6s = append(email.IPv6s, tmp)
					}

				}

			} else if IsIPv4(host) {

				ipv4 := ParseIPv4(host)

				if ipv4 != nil {

					tmp := ipv4.String()

					if !slices.Contains(email.IPv4s, tmp) {
						email.IPv4s = append(email.IPv4s, tmp)
					}

				}

			} else if IsDomain(host) {

				domain := ParseDomain(host)

				if domain != nil {

					tmp := domain.String()

					if !slices.Contains(email.Domains, tmp) {
						email.Domains = append(email.Domains, tmp)
					}

				}

			}

		}

	}

	// XXX: Must be done AFTER parsing Received path
	if raw, ok := headers["from"]; ok {

		tmp := strings.TrimSpace(raw)

		if strings.Contains(tmp, "<") && strings.HasSuffix(tmp, ">") {

			tmp_address := tmp[strings.Index(tmp, "<")+1:len(tmp)-1]

			if strings.Contains(tmp_address, "@") {

				address := strings.ToLower(strings.TrimSpace(tmp_address))

				email.From = address

				host := strings.ToLower(strings.TrimSpace(address[strings.Index(address, "@")+1:]))

				if IsDomain(host) {

					domain := ParseDomain(host)

					if domain != nil {

						tmp := domain.String()

						if !slices.Contains(email.Domains, tmp) {
							email.Domains = append(email.Domains, tmp)
						}

					}

				}

			}

		}

	}

	if raw, ok := headers["return-path"]; ok {

		address := strings.TrimSpace(raw)

		if strings.HasPrefix(address, "<") && strings.HasSuffix(address, ">") {
			address = address[1:len(address)-1]
		}

		if strings.Contains(address, "@") {

			host := strings.ToLower(strings.TrimSpace(address[strings.Index(address, "@")+1:]))
			domain := ParseDomain(host)

			if domain != nil {

				tmp := domain.String()

				if !slices.Contains(email.Domains, tmp) {
					email.Domains = append(email.Domains, tmp)
				}

			}

		}

	}

	if raw, ok := headers["subject"]; ok {

		tmp := strings.TrimSpace(raw)

		if tmp != "" {
			email.Subject = utils_strings.ToASCII(tmp)
		}

	}

	if raw, ok := headers["to"]; ok {

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

		}

	}

	if email.To == "" {

		if raw, ok := headers["delivered-to"]; ok {

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

			}

		}

	}

	if email.Boundary != "" {

		// E-Mail with MIME Boundaries
		raw_buffer := strings.TrimSpace(string(buffer))
		message := raw_buffer[strings.Index(raw_buffer, "\n--" + email.Boundary) + len(email.Boundary)+3:]
		message = strings.TrimSpace(message)
		message = strings.ReplaceAll(message, "\r\n", "\n")

		if strings.Contains(message, email.Boundary) {

			contents := strings.Split(message, "\n--" + email.Boundary + "\n")

			for _, content := range contents {

				if strings.Contains(content, "Content-Type: multipart/alternative") {

					// Outlook format, which uses nested MIME boundaries
					nested_boundary := parseBoundary(content)

					if nested_boundary != "" {

						nested_contents := strings.Split(message, "\n--" + nested_boundary + "\n")

						for _, nested_content := range nested_contents {

							if strings.Contains(nested_content, "Content-Type: text/plain") {

								if strings.Contains(nested_content, "Content-Transfer-Encoding: base64") {

									if strings.Contains(nested_content, "\n\n") {

										tmp1 := strings.TrimSpace(nested_content[strings.Index(nested_content, "\n\n"):])
										tmp2, err2 := base64.StdEncoding.DecodeString(tmp1)

										if err2 == nil {
											email.Message = toASCIIMessage(string(tmp2))
										}

									} else {
										// Invalid Boundary?
									}

								} else {

									if strings.Contains(nested_content, "\n\n") {
										email.Message = strings.TrimSpace(nested_content[strings.Index(nested_content, "\n\n"):])
									} else {
										email.Message = strings.TrimSpace(nested_content)
									}

								}

								break

							}

						}

					}

				} else if strings.Contains(content, "Content-Type: text/plain") {

					// Apple Mail / Thunderbird format
					if strings.Contains(content, "Content-Transfer-Encoding: base64") {

						if strings.Contains(content, "\n\n") {

							tmp1 := strings.TrimSpace(content[strings.Index(content, "\n\n"):])
							tmp2, err2 := base64.StdEncoding.DecodeString(tmp1)

							if err2 == nil {
								email.Message = toASCIIMessage(string(tmp2))
							}

						} else {
							// Invalid Boundary?
						}

					} else {

						if strings.Contains(content, "\n\n") {
							email.Message = strings.TrimSpace(content[strings.Index(content, "\n\n"):])
						} else {
							email.Message = strings.TrimSpace(content)
						}

					}

					break

				}

			}

		}

	} else {

		// E-Mail without MIME Boundaries
		raw_buffer := strings.TrimSpace(string(buffer))
		message := strings.TrimSpace(raw_buffer[strings.Index(raw_buffer, "\n\n\n") + 3:])

		if message != "" {
			email.Message = message
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
