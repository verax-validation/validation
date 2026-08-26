package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestRequiredString(t *testing.T) {
	if err := rules.Required("x"); err != nil {
		t.Errorf("Required(non-empty) = %v, want nil", err)
	}
	if err := rules.Required(""); err == nil {
		t.Error("Required(empty string) should fail")
	}
}

func TestRequiredScalars(t *testing.T) {
	cases := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"zero int", 0, true},
		{"non-zero int", 3, false},
		{"false bool", false, true},
		{"true bool", true, false},
		{"nil pointer", (*int)(nil), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.value.(type) {
			case int:
				err = rules.Required(v)
			case bool:
				err = rules.Required(v)
			case *int:
				err = rules.Required(v)
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("Required(%v) err = %v, wantErr %v", tc.value, err, tc.wantErr)
			}
		})
	}
}

func TestRequiredPointerToZeroValue(t *testing.T) {
	empty := ""
	if err := rules.Required(&empty); err == nil {
		t.Error("Required(pointer to empty string) should fail")
	}

	name := "alice"
	if err := rules.Required(&name); err != nil {
		t.Errorf("Required(pointer to non-empty) = %v, want nil", err)
	}
}

func TestOptionalSkipsZeroValue(t *testing.T) {
	executed := false
	alwaysFail := func(v int) error {
		executed = true
		return errFailInt
	}

	rule := rules.Optional(alwaysFail)

	if err := rule(0); err != nil {
		t.Errorf("Optional on zero value = %v, want nil", err)
	}
	if executed {
		t.Error("wrapped rule should not run for zero value")
	}
}

var errFailInt = errorf("fail")

func TestOptionalRunsNonZeroValue(t *testing.T) {
	rule := rules.Optional(rules.Max(100))

	if err := rule(150); err == nil {
		t.Error("Optional(Max(100))(150) should fail")
	}
	if err := rule(50); err != nil {
		t.Errorf("Optional(Max(100))(50) = %v, want nil", err)
	}
}

func TestRequiredIf(t *testing.T) {
	cond := false
	rule := rules.RequiredIf[string](func() bool { return cond })

	// when the condition is false, empty values pass
	cond = false
	if err := rule(""); err != nil {
		t.Errorf("RequiredIf(false)(empty) = %v, want nil", err)
	}

	// when the condition is true, empty values fail
	cond = true
	if err := rule(""); err == nil {
		t.Error("RequiredIf(true)(empty) should fail")
	}
	if err := rule("x"); err != nil {
		t.Errorf("RequiredIf(true)(x) = %v, want nil", err)
	}
}
