package verax_test

import (
	"testing"

	"github.com/verax-validation/validation"
)

type registerForm struct {
	Password string
	Confirm  string
	Start    int
	End      int
}

func TestFieldEq(t *testing.T) {
	f := &registerForm{Password: "abc", Confirm: "abc"}

	err := verax.ValidateStruct(f,
		verax.Field[string]().WithField(&f.Confirm).WithFieldEq(&f.Password),
	)
	if err != nil {
		t.Errorf("matching passwords = %v, want nil", err)
	}

	f.Confirm = "abd"
	err = verax.ValidateStruct(f,
		verax.Field[string]().WithField(&f.Confirm).WithFieldEq(&f.Password),
	)
	if err == nil {
		t.Fatal("mismatched passwords should fail")
	}
	if want := "confirm: must equal password"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestFieldNe(t *testing.T) {
	f := &registerForm{Password: "abc", Confirm: "abc"}

	err := verax.ValidateStruct(f,
		verax.Field[string]().WithField(&f.Confirm).WithFieldNe(&f.Password),
	)
	if err == nil {
		t.Fatal("confirm should differ from password")
	}
	if want := "confirm: must not equal password"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestFieldGte(t *testing.T) {
	f := &registerForm{Start: 5, End: 3}

	err := verax.ValidateStruct(f,
		verax.Field[int]().WithField(&f.End).WithFieldGte(&f.Start),
	)
	if err == nil {
		t.Fatal("end < start should fail")
	}
	if want := "end: must be at least start"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}

	f.End = 5
	if err := verax.ValidateStruct(f,
		verax.Field[int]().WithField(&f.End).WithFieldGte(&f.Start),
	); err != nil {
		t.Errorf("end == start should pass, got %v", err)
	}
}

func TestFieldLt(t *testing.T) {
	f := &registerForm{Start: 5, End: 5}

	err := verax.ValidateStruct(f,
		verax.Field[int]().WithField(&f.Start).WithFieldLt(&f.End),
	)
	if err == nil {
		t.Fatal("start == end should fail for strict Lt")
	}
	if want := "start: must be less than end"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}

	f.Start = 4
	if err := verax.ValidateStruct(f,
		verax.Field[int]().WithField(&f.Start).WithFieldLt(&f.End),
	); err != nil {
		t.Errorf("start < end should pass, got %v", err)
	}
}

func TestFieldCompareWithLabel(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	f := &registerForm{Password: "abc", Confirm: "abd"}

	err := verax.ValidateStruct(f,
		verax.Field[string]().
			WithField(&f.Confirm).
			WithLabel("确认密码").
			WithFieldEq(&f.Password),
	)
	if err == nil {
		t.Fatal("expected failure")
	}
	// the label prefix should apply to cross-field comparison errors
	if want := "confirm: 确认密码必须等于 password"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}
