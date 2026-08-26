package verax_test

import (
	"errors"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/rules"
)

type profile struct {
	Name string
	Age  int
}

func TestFieldOptionLabel(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	p := &profile{Name: ""}

	err := verax.ValidateStruct(p,
		verax.Field[string]().
			WithField(&p.Name).
			WithRules(rules.Required[string]).
			WithLabel("名称"),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	// label prefix concatenation: label + message
	want := "name: 名称不能为空"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFieldOptionWithoutLabel(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	p := &profile{Name: ""}

	err := verax.ValidateStruct(p,
		verax.Field[string]().WithField(&p.Name).WithRules(rules.Required[string]),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	// the original message is kept when no label is set
	want := "name: 不能为空"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFieldOptionCustomErr(t *testing.T) {
	p := &profile{Name: ""}

	err := verax.ValidateStruct(p,
		verax.Field[string]().
			WithField(&p.Name).
			WithRules(rules.Required[string]).
			WithErr("名称不能为空"),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	want := "name: 名称不能为空"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFieldOptionCheckFn(t *testing.T) {
	p := &profile{Name: "alice", Age: 15}

	// WithCheckFn is a parameterless check: business-level conditional validation
	err := verax.ValidateStruct(p,
		verax.Field[int]().
			WithField(&p.Age).
			WithCheckFn(func() error {
				if p.Age < 18 {
					return errors.New("must be an adult")
				}
				return nil
			}),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	want := "age: must be an adult"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	// passes when the condition is satisfied
	p.Age = 20
	if err := verax.ValidateStruct(p,
		verax.Field[int]().
			WithField(&p.Age).
			WithCheckFn(func() error {
				if p.Age < 18 {
					return errors.New("must be an adult")
				}
				return nil
			}),
	); err != nil {
		t.Errorf("adult age should pass, got %v", err)
	}
}

func TestFieldOptionLabelWithCheckFn(t *testing.T) {
	// WithLabel also applies to bare errors returned by WithCheckFn, preserving the original error chain
	p := &profile{Name: "alice", Age: 15}
	adultErr := errors.New("must be an adult")

	err := verax.ValidateStruct(p,
		verax.Field[int]().
			WithField(&p.Age).
			WithCheckFn(func() error {
				return adultErr
			}).
			WithLabel("年龄"),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	want := "age: 年龄must be an adult"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	// bare errors are wrapped by labelWrapped, errors.Is still matches the original error
	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	if !errors.Is(errs[0], adultErr) {
		t.Error("errors.Is should still match the inner error")
	}
}

func TestFieldOptionMixed(t *testing.T) {
	// all options combined: field + rules + label + custom error + custom check
	p := &profile{Name: "", Age: 15}

	err := verax.ValidateStruct(p,
		verax.Field[string]().
			WithField(&p.Name).
			WithRules(rules.Required[string]).
			WithLabel("名称").
			WithErr("请填写名称"),
		verax.Field[int]().
			WithField(&p.Age).
			WithRules(rules.Min(18)).
			WithCheckFn(func() error {
				if p.Age%2 != 0 {
					return errors.New("age must be even")
				}
				return nil
			}),
	)
	if err == nil {
		t.Fatal("expected failure")
	}

	// output follows the field declaration order
	want := "name: 请填写名称; age: age must be even"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
