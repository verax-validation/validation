package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestContains(t *testing.T) {
	rule := rules.Contains("@")

	if err := rule("a@b"); err != nil {
		t.Errorf("Contains(@)(a@b) = %v, want nil", err)
	}
	if err := rule("ab"); err == nil {
		t.Error("Contains(@)(ab) should fail")
	}
}

func TestStartWith(t *testing.T) {
	rule := rules.StartWith("http")

	if err := rule("https://a"); err != nil {
		t.Errorf("StartWith(http)(https://a) = %v, want nil", err)
	}
	if err := rule("ftp://a"); err == nil {
		t.Error("StartWith(http)(ftp://a) should fail")
	}
}

func TestEndWith(t *testing.T) {
	rule := rules.EndWith(".com")

	if err := rule("a.com"); err != nil {
		t.Errorf("EndWith(.com)(a.com) = %v, want nil", err)
	}
	if err := rule("a.org"); err == nil {
		t.Error("EndWith(.com)(a.org) should fail")
	}
}

func TestLen(t *testing.T) {
	rule := rules.Len[string](6)

	if err := rule("123456"); err != nil {
		t.Errorf("Len[string](6)(123456) = %v, want nil", err)
	}
	for _, v := range []string{"12345", "1234567"} {
		if err := rule(v); err == nil {
			t.Errorf("Len[string](6)(%q) should fail", v)
		}
	}
}

func TestStringRuleErrorMessages(t *testing.T) {
	if got := rules.Contains("@")("ab").Error(); got != "must contain @" {
		t.Errorf("Contains message = %q, want %q", got, "must contain @")
	}
	if got := rules.Len[string](6)("123").Error(); got != "the length must be exactly 6" {
		t.Errorf("Len message = %q, want %q", got, "the length must be exactly 6")
	}
}

func TestExcludes(t *testing.T) {
	rule := rules.Excludes("admin")

	if err := rule("user"); err != nil {
		t.Errorf("Excludes(admin)(user) = %v, want nil", err)
	}
	if err := rule("admin123"); err == nil {
		t.Error("Excludes(admin)(admin123) should fail")
	}
}

func TestContainsAny(t *testing.T) {
	rule := rules.ContainsAny("!@#")

	if err := rule("abc!def"); err != nil {
		t.Errorf("ContainsAny(!@#)(abc!def) = %v, want nil", err)
	}
	if err := rule("abcdef"); err == nil {
		t.Error("ContainsAny(!@#)(abcdef) should fail")
	}
}

func TestTrimSpace(t *testing.T) {
	rule := rules.TrimSpace(rules.Required)

	if err := rule("  alice  "); err != nil {
		t.Errorf("TrimSpace(Required)(  alice  ) = %v, want nil", err)
	}
	if err := rule("   "); err == nil {
		t.Error("TrimSpace(Required)(   ) should fail")
	}
}
