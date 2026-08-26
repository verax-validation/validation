package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestURL(t *testing.T) {
	valid := []string{
		"https://example.com",
		"http://example.com/path?q=1",
		"ftp://files.example.org:21/pub",
	}
	for _, v := range valid {
		if err := is.URL(v); err != nil {
			t.Errorf("URL(%q) = %v, want nil", v, err)
		}
	}

	invalid := []string{"", "example.com", "/relative/path", "https://"}
	for _, v := range invalid {
		if err := is.URL(v); err == nil {
			t.Errorf("URL(%q) should fail", v)
		}
	}
}

func TestIP(t *testing.T) {
	valid := []string{"192.168.0.1", "::1", "2001:db8::8a2e:370:7334", "::ffff:192.168.0.1"}
	for _, v := range valid {
		if err := is.IP(v); err != nil {
			t.Errorf("IP(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range []string{"", "256.1.1.1", "example.com"} {
		if err := is.IP(v); err == nil {
			t.Errorf("IP(%q) should fail", v)
		}
	}
}

func TestIPv4(t *testing.T) {
	if err := is.IPv4("10.0.0.255"); err != nil {
		t.Errorf("IPv4 = %v, want nil", err)
	}
	for _, v := range []string{"::1", "2001:db8::1"} {
		if err := is.IPv4(v); err == nil {
			t.Errorf("IPv4(%q) should fail", v)
		}
	}
}

func TestIPv6(t *testing.T) {
	valid := []string{"::1", "2001:db8::8a2e:370:7334", "::ffff:192.168.0.1"}
	for _, v := range valid {
		if err := is.IPv6(v); err != nil {
			t.Errorf("IPv6(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range []string{"192.168.0.1", ""} {
		if err := is.IPv6(v); err == nil {
			t.Errorf("IPv6(%q) should fail", v)
		}
	}
}

func TestMAC(t *testing.T) {
	valid := []string{
		"00:11:22:33:44:55",
		"00-11-22-33-44-55",
		"0011.2233.4455",
	}
	for _, v := range valid {
		if err := is.MAC(v); err != nil {
			t.Errorf("MAC(%q) = %v, want nil", v, err)
		}
	}
	if err := is.MAC("00:11:22"); err == nil {
		t.Error("MAC truncated should fail")
	}
}
