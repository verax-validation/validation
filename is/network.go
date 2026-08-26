package is

import (
	"net"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/verax-validation/validation"
)

// URL requires the string to be parseable as an absolute URL, i.e. must include a scheme and host,
// e.g. "https://example.com/a".
var URL verax.Rule[string] = func(value string) error {
	u, err := url.Parse(value)
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return ErrURL
	}
	return nil
}

// IP requires the string to be a valid IPv4 or IPv6 address.
var IP verax.Rule[string] = func(value string) error {
	if _, err := netip.ParseAddr(value); err != nil {
		return ErrIP
	}
	return nil
}

// IPv4 requires the string to be a valid dotted-decimal IPv4 address.
var IPv4 verax.Rule[string] = func(value string) error {
	addr, err := netip.ParseAddr(value)
	if err != nil || !addr.Is4() {
		return ErrIPv4
	}
	return nil
}

// IPv6 requires the string to be a valid IPv6 address.
var IPv6 verax.Rule[string] = func(value string) error {
	addr, err := netip.ParseAddr(value)
	if err != nil || !addr.Is6() {
		return ErrIPv6
	}
	return nil
}

// MAC requires the string to be a valid MAC address, supporting colon, hyphen, and dot separators.
var MAC verax.Rule[string] = func(value string) error {
	if _, err := net.ParseMAC(value); err != nil {
		return ErrMAC
	}
	return nil
}

// Subdomain requires the string to be a single valid subdomain label (1-63 characters, without dots).
var Subdomain verax.Rule[string] = func(value string) error {
	if subdomainPattern.MatchString(value) {
		return nil
	}
	return ErrSubdomain
}

// Domain requires the string to be a full domain (at least two labels, total length at most 255 characters).
var Domain verax.Rule[string] = func(value string) error {
	if len(value) <= 255 && domainPattern.MatchString(value) {
		return nil
	}
	return ErrDomain
}

// DNSName requires the string to be a valid DNS name (dot-separated labels, trailing root dot allowed, total length at most 253 characters).
var DNSName verax.Rule[string] = func(value string) error {
	if len(value) <= 253 && dnsNamePattern.MatchString(value) {
		return nil
	}
	return ErrDNSName
}

// Host requires the string to be a valid IP address or DNS name.
var Host verax.Rule[string] = func(value string) error {
	if _, err := netip.ParseAddr(value); err == nil {
		return nil
	}
	if _, err := netip.ParseAddrPort(value); err == nil {
		return ErrHost
	}
	if DNSName(value) == nil {
		return nil
	}
	return ErrHost
}

// Port requires the string to be a port in 1-65535.
var Port verax.Rule[string] = func(value string) error {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 || n > 65535 {
		return ErrPort
	}
	return nil
}

// E164 requires the string to be an E.164 international phone number (starts with +, at most 15 digits).
var E164 verax.Rule[string] = func(value string) error {
	if e164Pattern.MatchString(value) {
		return nil
	}
	return ErrE164
}

// DialString requires the string to be a "host:port" dial address,
// where host is a valid IP or DNS name and port is in 1-65535.
var DialString verax.Rule[string] = func(value string) error {
	host, portText, err := net.SplitHostPort(value)
	if err != nil {
		return ErrDialString
	}
	if Host(host) != nil || Port(portText) != nil {
		return ErrDialString
	}
	return nil
}

// CIDR requires the string to be a valid CIDR notation, e.g. "192.168.1.0/24" or "2001:db8::/32".
var CIDR verax.Rule[string] = func(value string) error {
	if _, _, err := net.ParseCIDR(value); err != nil {
		return ErrCIDR
	}
	return nil
}

// CIDRv4 requires the string to be a valid IPv4 CIDR notation.
var CIDRv4 verax.Rule[string] = func(value string) error {
	ip, _, err := net.ParseCIDR(value)
	if err != nil || ip.To4() == nil {
		return ErrCIDRv4
	}
	return nil
}

// CIDRv6 requires the string to be a valid IPv6 CIDR notation.
var CIDRv6 verax.Rule[string] = func(value string) error {
	ip, _, err := net.ParseCIDR(value)
	if err != nil || ip.To4() != nil {
		return ErrCIDRv6
	}
	return nil
}
