package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestLengthBoundaries(t *testing.T) {
	rule := rules.Length[string](2, 4)

	if err := rule("ab"); err != nil {
		t.Errorf("length 2 = %v, want nil", err)
	}
	if err := rule("abc"); err != nil {
		t.Errorf("length 3 = %v, want nil", err)
	}
	for _, s := range []string{"", "a", "abcde"} {
		if err := rule(s); err == nil {
			t.Errorf("Length(2,4)(%q) should fail", s)
		}
	}
}

func TestLengthBytes(t *testing.T) {
	rule := rules.Length[[]byte](1, 3)

	if err := rule([]byte("ab")); err != nil {
		t.Errorf("bytes length 2 = %v, want nil", err)
	}
	if err := rule([]byte("abcd")); err == nil {
		t.Error("bytes length 4 should fail")
	}
}

// Code validates a custom type whose underlying type is string
type Code string

func TestLengthNamedStringType(t *testing.T) {
	rule := rules.Length[Code](2, 4)

	if err := rule(Code("GO")); err != nil {
		t.Errorf("Code(\"GO\") = %v, want nil", err)
	}
	if err := rule(Code("G")); err == nil {
		t.Error("Code(\"G\") should fail")
	}
}

func TestLengthErrorMessage(t *testing.T) {
	err := rules.Length[string](5, 100)("abc")

	want := "the length must be between 5 and 100"
	if got := err.Error(); got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}
