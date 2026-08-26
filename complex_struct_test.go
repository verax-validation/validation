package verax_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	"github.com/verax-validation/validation/rules"
)

// ---- complex struct model ----
//
// covering typical real-business shapes:
//   - custom named type implementing Validatable (OrderID)
//   - nested struct implementing Validatable (Item)
//   - element-level aggregation of slice and map fields (collections)
//   - natural optional semantics of pointer fields (Note)
//   - multi-rule chains with short-circuit (Name)

type OrderID string

// Validate requires the order ID to be prefixed with ORD-
func (id OrderID) Validate() error {
	if len(id) == 0 || !strings.HasPrefix(string(id), "ORD-") {
		return errors.New("must be prefixed with ORD-")
	}
	return nil
}

type Item struct {
	Name string
	Qty  int
}

func (it Item) Validate() error {
	return verax.ValidateStruct(&it,
		verax.Field[string]().WithField(&it.Name).WithRules(rules.Required),
		verax.Field[int]().WithField(&it.Qty).WithRules(rules.Min(1)),
	)
}

type Shipment struct {
	ID     string
	Items  []Item
	Labels map[string]string
	Note   *string
	Rush   bool
	Reason string
}

var shipmentRules = map[string][]verax.Rule[string]{
	// a map literal lacks the target type context, so generic rules must be instantiated explicitly
	"id": {rules.Required[string], rules.Length[string](5, 20)},
}

func (s *Shipment) Validate() error {
	return verax.ValidateStruct(s,
		verax.Field[string]().WithField(&s.ID).WithRules(shipmentRules["id"]...),
		verax.Field[[]Item]().WithField(&s.Items).WithRules(collections.Slice[Item](verax.Valid)),
		verax.Field[map[string]string]().WithField(&s.Labels).WithRules(collections.Map[string, string](nonBlankLabel)),
		// rules for pointer fields need manual dereferencing: nil means unfilled and passes directly
		verax.Field[*string]().WithField(&s.Note).WithRules(func(v *string) error {
			if v == nil {
				return nil
			}
			return verax.Validate(*v, rules.Length[string](1, 200))
		}),
		// rush orders must provide a reason
		verax.Field[string]().WithField(&s.Reason).WithRules(verax.When[string](s.Rush, rules.Required)),
	)
}

func nonBlankLabel(v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return errors.New("label cannot be blank")
	}
	return nil
}

func validShipment() *Shipment {
	note := "handle with care"
	return &Shipment{
		ID:     "ORD-20260825-001",
		Items:  []Item{{Name: "keyboard", Qty: 2}, {Name: "mouse", Qty: 1}},
		Labels: map[string]string{"en": "fragile", "jp": "取扱注意"},
		Note:   &note,
	}
}

// TestComplexStructPass passes as a whole when all fields are valid
func TestComplexStructPass(t *testing.T) {
	if err := validShipment().Validate(); err != nil {
		t.Errorf("valid shipment should pass, got %v", err)
	}
}

// TestComplexStructAggregatedOutput when multiple levels fail at the same time,
// errors aggregate level by level along the path field -> index/key -> sub-field into readable text
func TestComplexStructAggregatedOutput(t *testing.T) {
	bad := validShipment()
	bad.ID = "X"            // length too short
	bad.Items[1].Name = ""  // Items.1.Name missing
	bad.Items[1].Qty = 0    // Items.1.Qty out of range
	bad.Labels["cn"] = "  " // Labels.cn blank
	bad.Note = new(string)  // Note explicitly set with an over-length value
	*bad.Note = strings.Repeat("x", 201)

	err := bad.Validate()
	if err == nil {
		t.Fatal("expected aggregated failure, got nil")
	}

	want := "id: the length must be between 5 and 20; " +
		"items: 1: name: cannot be blank; qty: must be no less than 1; " +
		"labels: cn: label cannot be blank; " +
		"note: the length must be between 1 and 200"
	if got := err.Error(); got != want {
		t.Errorf("aggregated output =\n%q\nwant\n%q", got, want)
	}
}

// TestComplexStructNilPointerSkipped nil pointer fields are treated as unfilled, so the whole still passes
func TestComplexStructNilPointerSkipped(t *testing.T) {
	s := validShipment()
	s.Note = nil

	if err := s.Validate(); err != nil {
		t.Errorf("nil optional pointer should pass, got %v", err)
	}
}

// TestComplexStructConditionalField fails for a rush order missing a reason, passes for a normal order without one
func TestComplexStructConditionalField(t *testing.T) {
	rush := validShipment()
	rush.Rush = true
	rush.Reason = ""

	err := rush.Validate()
	if err == nil {
		t.Fatal("rush order without reason should fail")
	}
	errs := verax.Errors{}
	if !errors.As(err, &errs) {
		t.Fatalf("type = %T, want verax.Errors", err)
	}
	if _, found := errs.Get("reason"); !found {
		t.Errorf("reason error missing in %v", errs)
	}

	normal := validShipment()
	normal.Rush = false
	normal.Reason = ""
	if err := normal.Validate(); err != nil {
		t.Errorf("normal order without reason should pass, got %v", err)
	}
}

// TestComplexStructLocalizedOutput switches the whole output language after injecting Chinese
func TestComplexStructLocalizedOutput(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	bad := validShipment()
	bad.Items[1].Name = ""

	want := "items: 1: name: 不能为空"
	if got := bad.Validate().Error(); got != want {
		t.Errorf("localized output =\n%q\nwant\n%q", got, want)
	}
}

// ---- simple value rule chains: short-circuit and placeholder ----

func requireAbc(v string) error {
	if !strings.Contains(v, "abc") {
		return errors.New("error abc")
	}
	return nil
}

func requireXyz(v string) error {
	if !strings.Contains(v, "xyz") {
		return errors.New("error xyz")
	}
	return nil
}

func TestRuleChainShortCircuit(t *testing.T) {
	cases := []struct {
		name  string
		value string
		rules []verax.Rule[string]
		want  string
	}{
		{"first fails", "123", []verax.Rule[string]{requireAbc, requireXyz}, "error abc"},
		{"second fails", "abc", []verax.Rule[string]{requireAbc, requireXyz}, "error xyz"},
		{"all pass", "abcxyz", []verax.Rule[string]{requireAbc, requireXyz}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := verax.Validate(tc.value, tc.rules...)
			switch tc.want {
			case "":
				if err != nil {
					t.Errorf("Validate() = %v, want nil", err)
				}
			default:
				if err == nil || err.Error() != tc.want {
					t.Errorf("Validate() = %v, want %q", err, tc.want)
				}
			}
		})
	}
}

func TestSkipInChainDoesNotBlock(t *testing.T) {
	// Skip always passes and does not block later rules
	if err := verax.Validate("abcxyz", requireAbc, verax.Skip[string], requireXyz); err != nil {
		t.Errorf("Skip should not block chain, got %v", err)
	}
	if err := verax.Validate("abc", requireAbc, verax.Skip[string], requireXyz); err == nil {
		t.Error("chain after Skip should still run and fail")
	}
}

// TestNamedTypeSliceMatchesOzzoScenario aligns with the element-wise aggregation scenario of named-type slices:
// valid elements interleaved with an invalid one are located by index
func TestNamedTypeSliceMatchesOzzoScenario(t *testing.T) {
	ids := []OrderID{"ORD-1", "", "ORD-2"}

	err := verax.Validate(ids, collections.Slice[OrderID](verax.Valid))
	if err == nil {
		t.Fatal("expected failure")
	}

	want := "1: must be prefixed with ORD-"
	if got := err.Error(); got != want {
		t.Errorf("slice output = %q, want %q", got, want)
	}
}
