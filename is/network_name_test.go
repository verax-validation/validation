package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestNetworkNameRules(t *testing.T) {
	checkRules(t, "Subdomain", is.Subdomain,
		[]string{"www", "api-gateway", "a"},
		[]string{"", "-abc", "abc-", "a.b", "a_b"})

	checkRules(t, "Domain", is.Domain,
		[]string{"example.com", "api.example.co.uk", "xn--bcher-kva.example"},
		[]string{"", "localhost", "example..com", "-a.example.com"})

	checkRules(t, "DNSName", is.DNSName,
		[]string{"example.com", "localhost", "a-b.c-d", "example.com."},
		[]string{"", "-bad.example.com", "bad-.example.com"})

	checkRules(t, "Host", is.Host,
		[]string{"192.168.0.1", "::1", "example.com"},
		[]string{"", "example..com"})
}

func TestPortRule(t *testing.T) {
	checkRules(t, "Port", is.Port,
		[]string{"1", "80", "65535"},
		[]string{"", "0", "65536", "80a"})
}

func TestE164(t *testing.T) {
	checkRules(t, "E164", is.E164,
		[]string{"+8613800138000", "+12345678901", "+1555555"},
		[]string{"", "13800138000", "+0123456789", "+", "+86 138 0013 8000"})
}

func TestDialString(t *testing.T) {
	checkRules(t, "DialString", is.DialString,
		[]string{"example.com:443", "127.0.0.1:8080", "[::1]:9000"},
		[]string{"", "example.com", ":8080", "example.com:", "example.com:0"})
}
