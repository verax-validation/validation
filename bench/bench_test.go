// Package bench benchmarks verax against ozzo-validation on the same scenarios.
// Three scenarios: simple value, multi-field struct, per-element slice, each with pass and fail paths.
package bench

import (
	"fmt"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	vis "github.com/verax-validation/validation/is"
	"github.com/verax-validation/validation/rules"
)

var (
	valuePass = "alice@example.com"
	valueFail = "hi"

	tags = makeTags(100)
)

func makeTags(n int) []string {
	items := make([]string, n)
	for i := range items {
		items[i] = fmt.Sprintf("tag-%03d", i)
	}
	return items
}

type user struct {
	Name  string
	Email string
	Age   int
}

func newUser() *user {
	return &user{Name: "alice", Email: "alice@example.com", Age: 30}
}

func brokenUser() *user {
	return &user{Name: "", Email: "not-an-email", Age: 15}
}

// ---- scenario 1: simple value ----

func BenchmarkVeraxValuePass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := verax.Validate(valuePass,
			rules.Required,
			rules.Length[string](5, 100),
			vis.EmailFormat,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOzzoValuePass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := validation.Validate(valuePass,
			validation.Required,
			validation.Length(5, 100),
			is.Email,
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVeraxValueFail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := verax.Validate(valueFail,
			rules.Required,
			rules.Length[string](5, 100),
		); err == nil {
			b.Fatal("expected failure")
		}
	}
}

func BenchmarkOzzoValueFail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := validation.Validate(valueFail,
			validation.Required,
			validation.Length(5, 100),
		); err == nil {
			b.Fatal("expected failure")
		}
	}
}

// ---- scenario 2: multi-field struct ----

func benchVeraxStruct(u *user) error {
	return verax.ValidateStruct(u,
		verax.Field[string]().WithField(&u.Name).WithRules(rules.Required, rules.Length[string](2, 64)),
		verax.Field[string]().WithField(&u.Email).WithRules(rules.Required, vis.EmailFormat),
		verax.Field[int]().WithField(&u.Age).WithRules(rules.Min(18), rules.Max(120)),
	)
}

func benchOzzoStruct(u *user) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(2, 64)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Age, validation.Min(18), validation.Max(120)),
	)
}

func BenchmarkVeraxStructPass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := benchVeraxStruct(newUser()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOzzoStructPass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := benchOzzoStruct(newUser()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVeraxStructFail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := benchVeraxStruct(brokenUser()); err == nil {
			b.Fatal("expected failure")
		}
	}
}

func BenchmarkOzzoStructFail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := benchOzzoStruct(brokenUser()); err == nil {
			b.Fatal("expected failure")
		}
	}
}

// ---- scenario 3: per-element slice ----

func BenchmarkVeraxSlicePass(b *testing.B) {
	rule := collections.Slice[string](
		rules.Required,
		rules.Length[string](1, 20),
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := rule(tags); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOzzoSlicePass(b *testing.B) {
	rule := validation.Each(
		validation.Required,
		validation.Length(1, 20),
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := rule.Validate(tags); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVeraxSliceFail(b *testing.B) {
	broken := append([]string{"go"}, tags...)
	broken[50] = ""
	rule := collections.Slice[string](rules.Required)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := rule(broken); err == nil {
			b.Fatal("expected failure")
		}
	}
}

func BenchmarkOzzoSliceFail(b *testing.B) {
	broken := append([]string{"go"}, tags...)
	broken[50] = ""
	rule := validation.Each(validation.Required)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := rule.Validate(broken); err == nil {
			b.Fatal("expected failure")
		}
	}
}
