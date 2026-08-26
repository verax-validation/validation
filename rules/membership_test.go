package rules_test

import (
	"strings"
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestInString(t *testing.T) {
	rule := rules.In("red", "green", "blue")

	for _, v := range []string{"red", "blue"} {
		if err := rule(v); err != nil {
			t.Errorf("In hit %q = %v, want nil", v, err)
		}
	}
	if err := rule("black"); err == nil {
		t.Error("In(\"black\") should fail")
	}
}

func TestInInt(t *testing.T) {
	rule := rules.In(1, 3, 5)

	if err := rule(3); err != nil {
		t.Errorf("In(1,3,5)(3) = %v, want nil", err)
	}
	if err := rule(2); err == nil {
		t.Error("In(1,3,5)(2) should fail")
	}
}

func TestNotIn(t *testing.T) {
	rule := rules.NotIn("admin", "root")

	if err := rule("user"); err != nil {
		t.Errorf("NotIn miss = %v, want nil", err)
	}
	if err := rule("admin"); err == nil {
		t.Error("NotIn hit should fail")
	}
}

func TestMatch(t *testing.T) {
	rule := rules.Match(`^[a-z]+$`)

	if err := rule("abc"); err != nil {
		t.Errorf("match = %v, want nil", err)
	}
	if err := rule("A1"); err == nil {
		t.Error("non-match should fail")
	}
}

func TestMatchInvalidPatternPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("invalid pattern should panic at construction time")
		}
	}()

	rules.Match("[invalid")
}

func TestDate(t *testing.T) {
	rule := rules.Date("2006-01-02")

	if err := rule("2026-08-25"); err != nil {
		t.Errorf("valid date = %v, want nil", err)
	}
	if err := rule("25/08/2026"); err == nil {
		t.Error("wrong format should fail")
	}
}

func TestDateErrorMessageContainsLayout(t *testing.T) {
	err := rules.Date("2006-01-02")("bad")

	if !strings.Contains(err.Error(), "2006-01-02") {
		t.Errorf("message %q should contain layout", err.Error())
	}
}
