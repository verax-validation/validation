package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestEq(t *testing.T) {
	rule := rules.Eq("admin")

	if err := rule("admin"); err != nil {
		t.Errorf("Eq(admin)(admin) = %v, want nil", err)
	}
	if err := rule("guest"); err == nil {
		t.Error("Eq(admin)(guest) should fail")
	}
}

func TestNe(t *testing.T) {
	rule := rules.Ne(0)

	if err := rule(1); err != nil {
		t.Errorf("Ne(0)(1) = %v, want nil", err)
	}
	if err := rule(0); err == nil {
		t.Error("Ne(0)(0) should fail")
	}
}

func TestGt(t *testing.T) {
	rule := rules.Gt(18)

	if err := rule(19); err != nil {
		t.Errorf("Gt(18)(19) = %v, want nil", err)
	}
	for _, v := range []int{18, 17} {
		if err := rule(v); err == nil {
			t.Errorf("Gt(18)(%d) should fail", v)
		}
	}
}

func TestLt(t *testing.T) {
	rule := rules.Lt(60)

	if err := rule(59); err != nil {
		t.Errorf("Lt(60)(59) = %v, want nil", err)
	}
	for _, v := range []int{60, 61} {
		if err := rule(v); err == nil {
			t.Errorf("Lt(60)(%d) should fail", v)
		}
	}
}

func TestCompareErrorMessages(t *testing.T) {
	if got := rules.Eq("admin")("guest").Error(); got != "must equal admin" {
		t.Errorf("Eq message = %q, want %q", got, "must equal admin")
	}
	if got := rules.Gt(18)(17).Error(); got != "must be greater than 18" {
		t.Errorf("Gt message = %q, want %q", got, "must be greater than 18")
	}
}
