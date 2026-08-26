package rules_test

import (
	"testing"

	"github.com/verax-validation/validation/rules"
)

func TestMin(t *testing.T) {
	rule := rules.Min(18)

	if err := rule(18); err != nil {
		t.Errorf("Min(18)(18) = %v, want nil", err)
	}
	if err := rule(17); err == nil {
		t.Error("Min(18)(17) should fail")
	}
}

func TestMinFloatInference(t *testing.T) {
	// T is inferred as float64 from the argument
	rule := rules.Min(0.5)

	if err := rule(0.6); err != nil {
		t.Errorf("Min(0.5)(0.6) = %v, want nil", err)
	}
	if err := rule(0.4); err == nil {
		t.Error("Min(0.5)(0.4) should fail")
	}
}

func TestMax(t *testing.T) {
	rule := rules.Max(100)

	if err := rule(100); err != nil {
		t.Errorf("Max(100)(100) = %v, want nil", err)
	}
	if err := rule(101); err == nil {
		t.Error("Max(100)(101) should fail")
	}
}

func TestBetweenBoundaries(t *testing.T) {
	rule := rules.Between(1, 10)

	for _, v := range []int{1, 5, 10} {
		if err := rule(v); err != nil {
			t.Errorf("Between(1,10)(%d) = %v, want nil", v, err)
		}
	}
	for _, v := range []int{0, 11} {
		if err := rule(v); err == nil {
			t.Errorf("Between(1,10)(%d) should fail", v)
		}
	}
}

// Score validates a custom type whose underlying type is float64
type Score float64

func TestNamedOrderedType(t *testing.T) {
	// T cannot be inferred as Score from 90, so instantiate explicitly
	rule := rules.Between[Score](0, 100)

	if err := rule(Score(88.5)); err != nil {
		t.Errorf("Between[Score](0,100)(88.5) = %v, want nil", err)
	}
	if err := rule(Score(-1)); err == nil {
		t.Error("Between[Score](0,100)(-1) should fail")
	}
}

func TestNumericErrorMessages(t *testing.T) {
	err := rules.Min(18)(17)
	want := "must be no less than 18"
	if got := err.Error(); got != want {
		t.Errorf("message = %q, want %q", got, want)
	}

	err = rules.Between(1, 10)(11)
	want = "must be between 1 and 10"
	if got := err.Error(); got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}
