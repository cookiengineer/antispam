package types

import "slices"
import "strings"

func between(value string, open string, close string) (string, bool) {

	start := strings.Index(value, open)

	if start < 0 {
		return "", false
	}

	start += len(open)
	end := strings.Index(value[start:], close)

	if end < 0 {
		return "", false
	}

	return value[start : start+end], true

}

func parseReceived(content string) []string {

	result := make([]string, 0)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	for _, raw := range lines {

		line := strings.TrimSpace(raw)

		if strings.HasPrefix(line, "from ") && strings.Contains(line, " by ") && strings.Contains(line, " with ") {

			by_index := strings.Index(line, " by ")
			with_index := strings.Index(line, " with ")

			if by_index <= 5 || with_index <= by_index {
				continue
			}

			from_chunk := line[5:by_index]
			by_chunk := line[by_index+4 : with_index]

			if strings.Contains(from_chunk, " (") && strings.Contains(from_chunk, ")") {

				from_domain := strings.TrimSpace(from_chunk[0:strings.Index(from_chunk, " (")])

				if IsDomain(from_domain) {

					if !slices.Contains(result, from_domain) {
						result = append(result, from_domain)
					}

				}

				from_details := ""

				if tmp, ok := between(from_chunk, " (", ")"); ok == true {
					from_details = strings.TrimSpace(tmp)
				}

				if strings.Contains(from_details, "[") && strings.Contains(from_details, "]") {

					from_ip := ""

					if tmp, ok := between(from_details, "[", "]"); ok == true {
						from_ip = strings.TrimSpace(tmp)
					}

					if strings.HasPrefix(from_ip, "IPV6:") {
						from_ip = strings.TrimSpace(from_ip[5:])
					}

					if IsIPv6(from_ip) {
						result = append(result, from_ip)
					} else if IsIPv4(from_ip) {
						result = append(result, from_ip)
					}

				} else if strings.Contains(from_details, " ") {

					from_domain := strings.TrimSpace(from_details[0:strings.Index(from_details, " ")])

					if IsDomain(from_domain) {

						if !slices.Contains(result, from_domain) {
							result = append(result, from_domain)
						}

					}

				} else {

					if IsDomain(from_details) {

						if !slices.Contains(result, from_details) {
							result = append(result, from_details)
						}

					}

				}

			} else if strings.Contains(from_chunk, " ") {

				from_domain := strings.TrimSpace(from_chunk[0:strings.Index(from_chunk, " ")])

				if IsDomain(from_domain) {

					if !slices.Contains(result, from_domain) {
						result = append(result, from_domain)
					}

				}

			}

			if strings.Contains(by_chunk, " (") && strings.Contains(by_chunk, ")") {

				by_domain := strings.TrimSpace(by_chunk[0:strings.Index(by_chunk, " (")])

				if IsDomain(by_domain) {

					if !slices.Contains(result, by_domain) {
						result = append(result, by_domain)
					}

				}

				by_details := ""

				if tmp, ok := between(by_chunk, " (", ")"); ok == true {
					by_details = strings.TrimSpace(tmp)
				}

				if strings.Contains(by_details, "[") && strings.Contains(by_details, "]") {

					by_ip := ""

					if tmp, ok := between(by_details, "[", "]"); ok == true {
						by_ip = strings.TrimSpace(tmp)
					}

					if strings.HasPrefix(by_ip, "IPV6:") {
						by_ip = strings.TrimSpace(by_ip[5:])
					}

					if IsIPv6(by_ip) {

						if !slices.Contains(result, by_ip) {
							result = append(result, by_ip)
						}

					} else if IsIPv4(by_ip) {

						if !slices.Contains(result, by_ip) {
							result = append(result, by_ip)
						}

					}

				} else if strings.Contains(by_details, " ") {

					by_domain := strings.TrimSpace(by_details[0:strings.Index(by_details, " ")])

					if IsDomain(by_domain) {

						if !slices.Contains(result, by_domain) {
							result = append(result, by_domain)
						}

					}

				} else {

					if IsDomain(by_details) {

						if !slices.Contains(result, by_details) {
							result = append(result, by_details)
						}

					}

				}

			} else if strings.Contains(by_chunk, " ") {

				by_domain := strings.TrimSpace(by_chunk[0:strings.Index(by_chunk, " ")])

				if IsDomain(by_domain) {

					if !slices.Contains(result, by_domain) {
						result = append(result, by_domain)
					}

				}

			} else {

				if IsDomain(by_chunk) {

					if !slices.Contains(result, by_chunk) {
						result = append(result, by_chunk)
					}

				}

			}

		} else if strings.HasPrefix(line, "by ") && strings.Contains(line, " with ") {

			with_index := strings.Index(line, " with ")

			if with_index <= 3 {
				continue
			}

			by_chunk := line[3:with_index]

			if strings.Contains(by_chunk, " (") && strings.Contains(by_chunk, ")") {

				by_domain := strings.TrimSpace(by_chunk[0:strings.Index(by_chunk, " (")])

				if IsDomain(by_domain) {

					if !slices.Contains(result, by_domain) {
						result = append(result, by_domain)
					}

				}

				by_details := ""

				if tmp, ok := between(by_chunk, " (", ")"); ok == true {
					by_details = strings.TrimSpace(tmp)
				}

				if strings.Contains(by_details, "[") && strings.Contains(by_details, "]") {

					by_ip := ""

					if tmp, ok := between(by_details, "[", "]"); ok == true {
						by_ip = strings.TrimSpace(tmp)
					}

					if strings.HasPrefix(by_ip, "IPV6:") {
						by_ip = strings.TrimSpace(by_ip[5:])
					}

					if IsIPv6(by_ip) {

						if !slices.Contains(result, by_ip) {
							result = append(result, by_ip)
						}

					} else if IsIPv4(by_ip) {

						if !slices.Contains(result, by_ip) {
							result = append(result, by_ip)
						}

					}

				} else if strings.Contains(by_details, " ") {

					by_domain := strings.TrimSpace(by_details[0:strings.Index(by_details, " ")])

					if IsDomain(by_domain) {

						if !slices.Contains(result, by_domain) {
							result = append(result, by_domain)
						}

					}

				} else {

					if IsDomain(by_details) {

						if !slices.Contains(result, by_details) {
							result = append(result, by_details)
						}

					}

				}

			} else if strings.Contains(by_chunk, " ") {

				by_domain := strings.TrimSpace(by_chunk[0:strings.Index(by_chunk, " ")])

				if IsDomain(by_domain) {

					if !slices.Contains(result, by_domain) {
						result = append(result, by_domain)
					}

				}

			} else {

				if IsDomain(by_chunk) {

					if !slices.Contains(result, by_chunk) {
						result = append(result, by_chunk)
					}

				}

			}

		} else {

			// TODO: Different Received header format!?

		}

	}

	return result

}
