package verax_test

import (
	"errors"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	"github.com/verax-validation/validation/rules"
)

// TestInlineGenericRulesInference verifies that bare generic rules inlined into Validate/Field need no explicit type parameters
func TestInlineGenericRulesInference(t *testing.T) {
	if err := verax.Validate("alice", rules.Required); err != nil {
		t.Errorf("Validate(name, rules.Required) = %v, want nil", err)
	}

	if err := verax.Validate("", rules.Required); err == nil {
		t.Error("Validate(empty, rules.Required) should fail")
	}
}

type Account struct {
	Name string
	Role string
}

func TestFieldWithBuiltinRules(t *testing.T) {
	account := &Account{Name: "alice", Role: "guest"}

	err := verax.ValidateStruct(account,
		verax.Field[string]().WithField(&account.Name).WithRules(rules.Required[string], rules.Length[string](2, 32)),
		verax.Field[string]().WithField(&account.Role).WithRules(rules.Required[string], rules.In("guest", "member")),
	)

	if err != nil {
		t.Errorf("ValidateStruct() = %v, want nil", err)
	}
}

func TestFieldWithBuiltinRulesFailure(t *testing.T) {
	account := &Account{Name: "", Role: "hacker"}

	err := verax.ValidateStruct(account,
		verax.Field[string]().WithField(&account.Name).WithRules(rules.Required[string]),
		verax.Field[string]().WithField(&account.Role).WithRules(rules.In("guest", "member")),
	)

	var errs verax.Errors
	ok := errors.As(err, &errs)
	if !ok {
		t.Fatalf("ValidateStruct() = %T, want verax.Errors", err)
	}
	if len(errs) != 2 {
		t.Errorf("error count = %d, want 2: %v", len(errs), errs)
	}
}

type Register struct {
	Password      string
	Confirm       string
	PaymentMethod string
	CardNumber    string
	Tags          []string
}

func TestNewRulesIntegration(t *testing.T) {
	form := &Register{
		Password:      "123456",
		Confirm:       "123456",
		PaymentMethod: "card",
		CardNumber:    "",
		Tags:          []string{"go", "go"},
	}

	err := verax.ValidateStruct(form,
		verax.Field[string]().WithField(&form.Confirm).WithFieldEq(&form.Password),
		verax.Field[string]().WithField(&form.CardNumber).WithRules(rules.RequiredIf[string](func() bool { return form.PaymentMethod == "card" })),
		verax.Field[[]string]().WithField(&form.Tags).WithRules(collections.Unique[string]()),
	)

	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("ValidateStruct() = %T, want verax.Errors", err)
	}
	// the conditional-required and collection-unique each contribute one failure, the cross-field comparison passes
	if len(errs) != 2 {
		t.Errorf("error count = %d, want 2: %v", len(errs), errs)
	}
	if _, found := errs.Get("card_number"); !found {
		t.Errorf("missing card_number failure in %v", errs)
	}
	if _, found := errs.Get("tags"); !found {
		t.Errorf("missing tags failure in %v", errs)
	}
}
