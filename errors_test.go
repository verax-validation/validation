package verax_test

import (
	"errors"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/is"
	"github.com/verax-validation/validation/rules"
)

func TestNewError(t *testing.T) {
	err := verax.NewError(rules.CodeRequired, "cannot be blank")

	if err.Code != rules.CodeRequired {
		t.Errorf("Code = %q, want %q", err.Code, rules.CodeRequired)
	}
	// Message is auto-filled by the constructor from the built-in English template table
	if err.Message != "cannot be blank" {
		t.Errorf("Message = %q, want %q", err.Message, "cannot be blank")
	}
	if got := err.Error(); got != "cannot be blank" {
		t.Errorf("Error() = %q, want %q", got, "cannot be blank")
	}
}

func TestNewMessageRender(t *testing.T) {
	// parameters are interpolated via Go template syntax
	err := verax.NewMessage(rules.CodeLength, map[string]string{
		"min": "5",
		"max": "100",
	})

	if err.Code != rules.CodeLength {
		t.Errorf("Code = %q, want %q", err.Code, rules.CodeLength)
	}
	want := "the length must be between 5 and 100"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestNewMessageRendersCurrentLocale(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	want := "不能小于 18"
	if got := verax.NewMessage(rules.CodeMin, map[string]string{"min": "18"}).Error(); got != want {
		t.Errorf("zh render = %q, want %q", got, want)
	}
}

func TestUnknownCodeFallsBackToCodeItself(t *testing.T) {
	// with an unknown code and no text, the code itself is rendered for diagnosis
	if got := verax.NewMessage("app.undefined.code", nil).Error(); got != "app.undefined.code" {
		t.Errorf("Error() = %q, want the code itself", got)
	}
}

func TestCustomCodeMessage(t *testing.T) {
	// application custom error: the full message is constructed directly
	err := verax.NewError("app.token.expired", "登录已过期")

	if err.Code != "app.token.expired" {
		t.Errorf("Code = %q, want app.token.expired", err.Code)
	}
	if got := err.Error(); got != "登录已过期" {
		t.Errorf("Error() = %q, want %q", got, "登录已过期")
	}
}

func TestErrorsFormat(t *testing.T) {
	// bare errors without field attribution are located by validation-order index
	errs := verax.Errors{
		errors.New("cannot be blank"),
		errors.New("must be a valid email address"),
		errors.New("must be no less than 1"),
	}

	want := "0: cannot be blank; 1: must be a valid email address; 2: must be no less than 1"
	if got := errs.Error(); got != want {
		t.Errorf("Errors() =\n%q\nwant\n%q", got, want)
	}
}

func TestErrorsEmpty(t *testing.T) {
	var errs verax.Errors

	if got := errs.Error(); got != "" {
		t.Errorf("empty Errors = %q, want empty string", got)
	}
}

func TestSentinelMessagesNotEmpty(t *testing.T) {
	// sentinel errors of built-in rules must have usable English text;
	// new rules missing text are reported here
	sentinels := map[string]*verax.Error{
		"ErrStructNil":      verax.ErrStructNil,
		"rules.ErrRequired": rules.ErrRequired,
		"rules.ErrIn":       rules.ErrIn,
		"rules.ErrNotIn":    rules.ErrNotIn,
		"rules.ErrMatch":    rules.ErrMatch,
		"is.ErrEmail":       is.ErrEmail,
		"is.ErrURL":         is.ErrURL,
		"is.ErrUUID":        is.ErrUUID,
		"is.ErrCreditCard":  is.ErrCreditCard,
	}
	for name, e := range sentinels {
		if len(e.Message) == 0 {
			t.Errorf("%s has empty Message", name)
		}
	}
}
