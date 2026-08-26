package verax_test

import (
	"testing"

	verax "github.com/verax-validation/validation"
	"github.com/verax-validation/validation/rules"
)

func TestRegisterBuiltinLocales(t *testing.T) {
	// each Render produces a new error in the currently active language
	rule := rules.Required[string]

	cases := []struct {
		register func()
		locale   string
		want     string
	}{
		{verax.RegisterEn, verax.LocaleEn, "cannot be blank"},
		{verax.RegisterZhCN, verax.LocaleZhCN, "不能为空"},
		{verax.RegisterZhTW, verax.LocaleZhTW, "不能為空"},
		{verax.RegisterJa, verax.LocaleJa, "必須項目です"},
		{verax.RegisterFr, verax.LocaleFr, "ne peut pas être vide"},
		{verax.RegisterDe, verax.LocaleDe, "darf nicht leer sein"},
	}
	for _, tc := range cases {
		tc.register()

		if got := rule("").Error(); got != tc.want {
			t.Errorf("locale %s: Error() = %q, want %q", tc.locale, got, tc.want)
		}
	}
	verax.RegisterEn()
}
