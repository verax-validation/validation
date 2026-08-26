package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestCIDR(t *testing.T) {
	checkRules(t, "CIDR", is.CIDR,
		[]string{"192.168.1.0/24", "10.0.0.0/8", "2001:db8::/32", "0.0.0.0/0"},
		[]string{"", "192.168.1.0", "192.168.1.0/33", "2001:db8::/129"})
}

func TestCIDRv4(t *testing.T) {
	checkRules(t, "CIDRv4", is.CIDRv4,
		[]string{"192.168.1.0/24", "0.0.0.0/0"},
		[]string{"", "2001:db8::/32", "192.168.1.0/33"})
}

func TestCIDRv6(t *testing.T) {
	checkRules(t, "CIDRv6", is.CIDRv6,
		[]string{"2001:db8::/32", "::/0"},
		[]string{"", "192.168.1.0/24", "2001:db8::/129"})
}
