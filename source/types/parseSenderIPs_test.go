package types

import "testing"

func TestParseSenderIPsHeadersOnly(t *testing.T) {

	mail := "Return-Path: <sender@example.com>\r\n" +
		"Received-SPF: pass (example.com: domain of sender@example.com designates 203.0.113.5 as permitted sender) client-ip=203.0.113.5;\r\n" +
		"X-Originating-IP: [198.51.100.7]\r\n" +
		"\r\n" +
		"Body text mentions designates 192.0.2.99 as permitted sender\r\n"

	ips := parseSenderIPs(headerBlock([]byte(mail)))

	if len(ips) != 2 {
		t.Fatalf("Expected 2 sender IPs, got %v", ips)
	}

	found_spf := false
	found_origin := false

	for _, ip := range ips {

		if ip == "203.0.113.5" {
			found_spf = true
		}

		if ip == "198.51.100.7" {
			found_origin = true
		}

		if ip == "192.0.2.99" {
			t.Errorf("Body IP should not be parsed as sender IP")
		}

	}

	if found_spf == false {
		t.Errorf("Expected SPF designated IP")
	}

	if found_origin == false {
		t.Errorf("Expected X-Originating-IP")
	}

}
