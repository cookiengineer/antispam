package types

import "bytes"
import "regexp"
import "strings"

var sender_ip_patterns []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`(?i)designates\s+\[?([0-9a-f:.]+)\]?\s+as\s+permitted\s+sender`),
	regexp.MustCompile(`(?i)client-ip=([0-9a-f:.]+)`),
	regexp.MustCompile(`(?i)smtp\.remote-ip=([0-9a-f:.]+)`),
	regexp.MustCompile(`(?i)x-originating-ip:\s*\[?([0-9a-f:.]+)\]?`),
}

func appendSenderIP(result []string, value string) []string {

	value = strings.TrimSpace(value)
	value = strings.Trim(value, "[]")
	value = strings.TrimRight(value, ".,;")

	if IsIPv6(value) == true {

		ipv6 := ParseIPv6(value)

		if ipv6 != nil {

			tmp := ipv6.String()

			for _, item := range result {

				if item == tmp {
					return result
				}

			}

			return append(result, tmp)

		}

	} else if IsIPv4(value) == true {

		ipv4 := ParseIPv4(value)

		if ipv4 != nil {

			tmp := ipv4.String()

			for _, item := range result {

				if item == tmp {
					return result
				}

			}

			return append(result, tmp)

		}

	}

	return result

}

func headerBlock(buffer []byte) []byte {

	limit := 256 * 1024

	if index := bytes.Index(buffer, []byte("\r\n\r\n")); index >= 0 {

		if index > limit {
			index = limit
		}

		return buffer[:index]

	}

	if index := bytes.Index(buffer, []byte("\n\n")); index >= 0 {

		if index > limit {
			index = limit
		}

		return buffer[:index]

	}

	if len(buffer) > limit {
		return buffer[:limit]
	}

	return buffer

}

func parseSenderIPs(buffer []byte) []string {

	result := make([]string, 0)
	content := string(buffer)

	for _, pattern := range sender_ip_patterns {

		matches := pattern.FindAllStringSubmatch(content, -1)

		for _, match := range matches {

			if len(match) >= 2 {
				result = appendSenderIP(result, match[1])
			}

		}

	}

	return result

}
