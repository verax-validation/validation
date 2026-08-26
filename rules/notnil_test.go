package rules_test

import (
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/rules"
)

func TestNotNil(t *testing.T) {
	var p *int
	if err := verax.Validate(p, rules.NotNil); err == nil {
		t.Error("NotNil(nil) should fail")
	}

	v := 3
	if err := verax.Validate(&v, rules.NotNil); err != nil {
		t.Errorf("NotNil(&v) = %v, want nil", err)
	}
}

func TestNotNilMessage(t *testing.T) {
	var p *string
	if got := verax.Validate(p, rules.NotNil).Error(); got != "must not be nil" {
		t.Errorf("message = %q, want %q", got, "must not be nil")
	}
}
