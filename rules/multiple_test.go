package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestMultipleOf(t *testing.T) {
	rule := rules.MultipleOf(10)

	for _, v := range []int{10, 20, 100} {
		if err := rule(v); err != nil {
			t.Errorf("MultipleOf(10)(%d) = %v, want nil", v, err)
		}
	}
	for _, v := range []int{5, 15} {
		if err := rule(v); err == nil {
			t.Errorf("MultipleOf(10)(%d) should fail", v)
		}
	}
}

func TestMultipleOfZeroStep(t *testing.T) {
	// with a zero base nothing is divisible, so it always fails
	if err := rules.MultipleOf(0)(1); err == nil {
		t.Error("MultipleOf(0)(1) should fail")
	}
}

func TestMultipleOfMessage(t *testing.T) {
	if got := rules.MultipleOf(10)(5).Error(); got != "must be a multiple of 10" {
		t.Errorf("message = %q, want %q", got, "must be a multiple of 10")
	}
}
