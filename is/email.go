package is

import (
	"net"
	"strings"

	"github.com/verax-validation/validation"
)

// Email requires the string to be a well-formed email whose domain has MX records (requires a DNS lookup).
// Use EmailFormat in network-restricted environments.
var Email verax.Rule[string] = func(value string) error {
	if err := EmailFormat(value); err != nil {
		return err
	}
	domain := value[strings.LastIndex(value, "@")+1:]
	records, err := net.LookupMX(domain)
	if err != nil || len(records) == 0 {
		return ErrEmail
	}
	return nil
}

// EmailFormat only requires the string to be a well-formed email, without any DNS lookup.
var EmailFormat verax.Rule[string] = func(value string) error {
	if emailPattern.MatchString(value) {
		return nil
	}
	return ErrEmail
}
