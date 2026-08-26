package is_test

import (
	"testing"

	"github.com/verax-validation/validation/is"
)

// checkRules uniformly runs the valid and invalid case tables
func checkRules(t *testing.T, name string, rule func(string) error, valid, invalid []string) {
	t.Helper()
	for _, v := range valid {
		if err := rule(v); err != nil {
			t.Errorf("%s(%q) = %v, want nil", name, v, err)
		}
	}
	for _, v := range invalid {
		if err := rule(v); err == nil {
			t.Errorf("%s(%q) should fail", name, v)
		}
	}
}

func TestCharsetRules(t *testing.T) {
	cases := []struct {
		name    string
		rule    func(string) error
		valid   []string
		invalid []string
	}{
		{"Alpha", is.Alpha,
			[]string{"abc", "ABC"},
			[]string{"", "ab1", "中文"}},
		{"Alphanumeric", is.Alphanumeric,
			[]string{"abc123"},
			[]string{"", "a-b", "a b"}},
		{"Digit", is.Digit,
			[]string{"0123456789"},
			[]string{"", "12a", "1.5"}},
		{"UTFLetter", is.UTFLetter,
			[]string{"abc", "中文", "日本語", "ñandú"},
			[]string{"", "ab1", "中1文"}},
		{"UTFDigit", is.UTFDigit,
			[]string{"0123", "١٢٣"},
			[]string{"", "12a", "一二三"}},
		{"UTFNumeric", is.UTFNumeric,
			[]string{"0123", "Ⅷ", "½"},
			[]string{"", "12a"}},
		{"UTFLetterNumeric", is.UTFLetterNumeric,
			[]string{"abc123", "中文Ⅷ"},
			[]string{"", "a-b", "!@#"}},
		{"LowerCase", is.LowerCase,
			[]string{"abc", "abc123"},
			[]string{"", "ABC", "Abc"}},
		{"UpperCase", is.UpperCase,
			[]string{"ABC", "ABC123"},
			[]string{"", "abc", "Abc"}},
		{"ASCII", is.ASCII,
			[]string{"abc 123!@#"},
			[]string{"", "中文", "café"}},
		{"PrintableASCII", is.PrintableASCII,
			[]string{"abc 123!@~"},
			[]string{"", "line\nbreak", "\x00"}},
		{"Hexadecimal", is.Hexadecimal,
			[]string{"deadBEEF", "1234567890"},
			[]string{"", "0x1f", "xyz"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			checkRules(t, tc.name, tc.rule, tc.valid, tc.invalid)
		})
	}
}
