package verax

import (
	"testing"

	"github.com/verax-validation/validation/internal/messages"
)

// TestBuiltinLocalesCompleteness verifies that the built-in language tables cover all error code constants.
// The tables are maintained in the same repo as the codes constants; new rules missing translations are reported here,
// avoiding a silent runtime fallback to English.
func TestBuiltinLocalesCompleteness(t *testing.T) {
	tables := map[string]map[string]string{
		"en":    messages.En,
		"zh-CN": messages.ZhCN,
		"zh-TW": messages.ZhTW,
		"ja":    messages.Ja,
		"fr":    messages.Fr,
		"de":    messages.De,
	}

	for name, table := range tables {
		// the English table is the text baseline; other languages must cover it code by code and introduce no extra entries
		for code := range messages.En {
			if _, ok := table[code]; !ok {
				t.Errorf("%s missing translation for %q", name, code)
			}
		}
		if got, want := len(table), len(messages.En); got != want {
			t.Errorf("%s has %d entries, want %d (no extra keys)", name, got, want)
		}
	}
}
