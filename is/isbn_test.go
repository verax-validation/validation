package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

func TestISBNRules(t *testing.T) {
	checkRules(t, "ISBN10", is.ISBN10,
		[]string{"0-306-40615-2", "0 306 40615 2"},
		[]string{"", "0306406151", "0306406153"})

	checkRules(t, "ISBN13", is.ISBN13,
		[]string{"978-3-16-148410-0", "978 3 16 148410 0"},
		[]string{"", "978-3-16-148410-1"})

	checkRules(t, "ISBN", is.ISBN,
		[]string{"0-306-40615-2", "978-3-16-148410-0"},
		[]string{"", "12345"})
}
