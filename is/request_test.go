package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestRequestRules(t *testing.T) {
	checkRules(t, "RequestURL", is.RequestURL,
		[]string{"https://example.com", "http://example.com/path?q=1"},
		[]string{"", "/relative", "example.com"})

	checkRules(t, "RequestURI", is.RequestURI,
		[]string{"/api/v1/users?q=1", "https://example.com/a", "/index.html"},
		[]string{"", "::bad-uri", "plain-text"})
}
