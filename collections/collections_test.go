package collections_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	"github.com/verax-validation/validation/rules"
)

var errBlank = errors.New("cannot be blank")

// mustGet extracts the error of the given field from the aggregate, failing directly when missing
func mustGet(t *testing.T, errs verax.Errors, field string) error {
	t.Helper()
	err, found := errs.Get(field)
	if !found {
		t.Fatalf("missing field %q in %v", field, errs)
	}
	return err
}

func nonBlank(v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return errBlank
	}
	return nil
}

func TestSliceAllPassed(t *testing.T) {
	rule := collections.Slice(nonBlank)

	if err := rule([]string{"a", "b"}); err != nil {
		t.Errorf("Slice = %v, want nil", err)
	}
}

func TestSliceAggregatesByIndex(t *testing.T) {
	rule := collections.Slice(nonBlank)

	err := rule([]string{"a", "", "c", ""})
	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	for _, want := range []string{"1", "3"} {
		if _, found := errs.Get(want); !found {
			t.Errorf("missing index %q in %v", want, errs)
		}
	}
	if got := len(errs); got != 2 {
		t.Errorf("error count = %d, want 2", got)
	}
	e1, found := errs.Get("1")
	if !found || !errors.Is(e1, errBlank) {
		t.Errorf("errs[1] = %v, want errBlank", e1)
	}
}

func TestSliceEmptyAndNil(t *testing.T) {
	rule := collections.Slice[string](nonBlank)

	if err := rule(nil); err != nil {
		t.Errorf("nil slice should pass, got %v", err)
	}
	if err := rule([]string{}); err != nil {
		t.Errorf("empty slice should pass, got %v", err)
	}
}

func TestSliceShortCircuitPerElement(t *testing.T) {
	calls := 0
	counting := func(v int) error {
		calls++
		return nil
	}
	failing := func(v int) error { return errBlank }

	// per-element short-circuit: the second rule only runs when the element passes the first
	rule := collections.Slice(failing, counting)
	_ = rule([]int{1, 2})

	if calls != 0 {
		t.Errorf("second rule ran %d times, want 0", calls)
	}
}

func TestEachOverSeq(t *testing.T) {
	rule := collections.Each(nonBlank)

	seq := slices.Values([]string{"x", "", "z"})
	err := rule(seq)

	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	if _, found := errs.Get("1"); !found {
		t.Errorf("missing index 1 in %v", errs)
	}
}

func TestMapValidatesValues(t *testing.T) {
	rule := collections.Map[string, string](nonBlank)

	err := rule(map[string]string{"a": "ok", "b": ""})
	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	if _, found := errs.Get("b"); !found {
		t.Errorf("missing key \"b\" in %v", errs)
	}
}

func TestLenRules(t *testing.T) {
	sliceRule := collections.SliceLen[int](1, 3)

	for _, items := range [][]int{{1}, {1, 2, 3}} {
		if err := sliceRule(items); err != nil {
			t.Errorf("SliceLen(1,3)(%v) = %v, want nil", items, err)
		}
	}
	for _, items := range [][]int{nil, {1, 2, 3, 4}} {
		if err := sliceRule(items); err == nil {
			t.Errorf("SliceLen(1,3)(%v) should fail", items)
		}
	}

	mapRule := collections.MapLen[string, int](2, 2)
	if err := mapRule(map[string]int{"a": 1, "b": 2}); err != nil {
		t.Errorf("MapLen(2,2) = %v, want nil", err)
	}
	if err := mapRule(map[string]int{"a": 1}); err == nil {
		t.Error("MapLen(2,2) with 1 entry should fail")
	}
}

func TestLenErrorMessage(t *testing.T) {
	err := collections.SliceLen[int](1, 3)(nil)

	want := "the length must be between 1 and 3"
	if got := err.Error(); got != want {
		t.Errorf("message = %q, want %q", got, want)
	}
}

func TestUnique(t *testing.T) {
	rule := collections.Unique[string]()

	if err := rule([]string{"a", "b", "c"}); err != nil {
		t.Errorf("Unique on distinct items = %v, want nil", err)
	}
	if err := rule([]string{"a", "b", "a"}); err == nil {
		t.Error("Unique on duplicated items should fail")
	}
	if err := rule(nil); err != nil {
		t.Errorf("Unique(nil) = %v, want nil", err)
	}
}

// Tagged is a nested struct whose collection elements participate via Valid
type Tagged struct {
	Name string
}

func (t Tagged) Validate() error {
	return verax.Validate(t.Name, rules.Required)
}

func TestSliceOfValidatable(t *testing.T) {
	form := struct {
		Tags []Tagged
	}{Tags: []Tagged{{Name: "go"}, {}}}

	err := verax.ValidateStruct(&form,
		// Valid is a generic function whose element type cannot be inferred for collections; instantiate it explicitly by convention
		verax.Field[[]Tagged]().WithField(&form.Tags).WithRules(collections.Slice[Tagged](verax.Valid)),
	)

	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	tagErrs, ok := errors.AsType[verax.Errors](mustGet(t, errs, "tags"))
	if !ok {
		t.Fatalf("errs[tags] type = %T, want verax.Errors", mustGet(t, errs, "tags"))
	}
	if _, found := tagErrs.Get("1"); !found {
		t.Errorf("nested failure should be at tags.1, got %v", tagErrs)
	}
}
