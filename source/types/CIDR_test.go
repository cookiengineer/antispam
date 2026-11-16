package types

import "testing"

func TestParseNetworks(t *testing.T) {

	tmp := ParseNetworks("192.168.0.0/24")

	if len(tmp) != 1 || tmp[0].String() != "192.168.0.0/24" {
		t.Errorf("Expected /24, got %v", tmp)
	}

	tmp = ParseNetworks("192.168.0.1 - 192.168.0.6")

	if len(tmp) != 4 {
		t.Errorf("Expected 4 prefixes, got %d (%v)", len(tmp), tmp)
	}

	if NormalizeNetwork("54.161.105.65/12") != "54.160.0.0/12" {
		t.Errorf("Expected network address to be normalized")
	}

	tmp = ParseNetworks("159.112.240.0 - 159.112.255.255")

	if len(tmp) != 1 || tmp[0].String() != "159.112.240.0/20" {
		t.Errorf("Expected /20 range conversion, got %v", tmp)
	}

	if NetworkContains("192.168.0.0/24", "192.168.0.42") == false {
		t.Errorf("Expected IP to be contained")
	}

	if NetworkContains("192.168.0.0/24", "192.168.1.42") == true {
		t.Errorf("Expected IP to be outside")
	}

	if ObservedPrefix("192.168.1.42") != "192.168.1.0/24" {
		t.Errorf("Expected observed /24")
	}

	if ObservedPrefix("2a01:e5c0:1234:5678::1") != "2a01:e5c0:1234::/48" {
		t.Errorf("Expected observed /48")
	}

	if IsPrivateIP("192.168.1.1") == false {
		t.Errorf("Expected private IP")
	}

	if IsPrivateIP("8.8.8.8") == true {
		t.Errorf("Expected public IP")
	}

}
