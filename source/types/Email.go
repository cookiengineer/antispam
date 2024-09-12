package types

import utils_strings "antispam/utils/strings"
import "slices"
import "strings"
import "time"

func parseReceivedLine(email *Email, line string) {

	if strings.Contains(line, " by ") {
		line = strings.TrimSpace(line[0:strings.Index(line, " by ")])
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

	received_multiline := ""
	section := ""

	for l := 0; l < len(lines); l++ {

		line := lines[l]

		if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {

			if section == "Received" {
				parseReceivedLine(&email, received_multiline)
				received_multiline = ""
			}

			section = strings.TrimSpace(line[0:strings.Index(line, ":")])

		}

		if strings.HasPrefix(line, "Return-Path: ") {

			address := strings.TrimSpace(line[13:])

			if strings.HasPrefix(address, "<") && strings.HasSuffix(address, ">") {
				address = address[1:len(address)-1]
			}

			if strings.Contains(address, "@") {

				domain := strings.TrimSpace(address[strings.Index(address, "@")+1:])

				if domain != "" {
					email.Domains = append(email.Domains, domain)
				}

			}

		} else if strings.HasPrefix(line, "Content-Type: ") || section == "Content-Type" {

			if strings.HasPrefix(line, "Content-Type: ") {

				tmp1 := strings.TrimSpace(line[14:])

				if strings.HasSuffix(tmp1, ";") {
					// Do Nothing
				} else if strings.Contains(tmp1, ";") {

					tmp2 := strings.Split(tmp1, ";")

					for t := 0; t < len(tmp2); t++ {

						chunk := strings.TrimSpace(tmp2[t])

						if strings.HasPrefix(chunk, "boundary=") {
							email.Boundary = utils_strings.TrimQuotes(strings.TrimSpace(chunk[9:]))
						}

					}

				} else {
					// Do Nothing
				}

			} else {

				chunk := strings.TrimSpace(line)

				if strings.HasPrefix(chunk, "boundary=") {
					email.Boundary = utils_strings.TrimQuotes(strings.TrimSpace(chunk[9:]))
				}

			}

		} else if strings.HasPrefix(line, "Message-ID: ") {

			tmp1 := strings.TrimSpace(line[12:])

			if strings.HasPrefix(tmp1, "<") && strings.HasSuffix(tmp1, ">") {
				email.MessageID = tmp1[1:len(tmp1)-1]
			}

		} else if strings.HasPrefix(line, "From: ") {

			tmp1 := strings.TrimSpace(line[6:])

			if strings.Contains(tmp1, "<") && strings.HasSuffix(tmp1, ">") {
				tmp1 = tmp1[1:len(tmp1)-1]
			}

			address := strings.TrimSpace(tmp1)
			domain := strings.TrimSpace(address[strings.Index(address, "@")+1:])

			email.From = address

			if !slices.Contains(email.Domains, domain) {
				email.Domains = append(email.Domains, domain)
			}

		} else if strings.HasPrefix(line, "To: ") {

			tmp1 := strings.TrimSpace(line[4:])

			if strings.Contains(tmp1, "<") && strings.HasSuffix(tmp1, ">") {
				tmp1 = tmp1[1:len(tmp1)-1]
			}

			address := strings.TrimSpace(tmp1)
			domain := strings.TrimSpace(address[strings.Index(address, "@")+1:])

			email.To = address

			if !slices.Contains(email.Domains, domain) {
				email.Domains = append(email.Domains, domain)
			}

		} else if strings.HasPrefix(line, "Subject: ") {

			tmp1 := strings.TrimSpace(line[9:])

			if tmp1 != "" {
				email.Subject = utils_strings.ToASCII(tmp1)
			}

		} else if strings.HasPrefix(line, "Received: ") || section == "Received" {

			if strings.HasPrefix(line, "Received: ") {
				received_multiline = strings.TrimSpace(line[10:])
			} else {
				received_multiline += " " + strings.TrimSpace(line)
			}

		} else {

			section = ""

		}

	}

	if email.Boundary != "" {

		raw_buffer := strings.TrimSpace(string(buffer))
		message := strings.TrimSpace(raw_buffer[strings.Index(raw_buffer, "--" + email.Boundary) + len(email.Boundary)+2:])

		if strings.Contains(message, email.Boundary) {
			message = message[0:strings.Index(message, "--" + email.Boundary)]
		}

		email.Message = strings.TrimSpace(message)

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
