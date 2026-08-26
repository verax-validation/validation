package verax

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"

	"github.com/verax-validation/validation/internal/codes"
)

// ErrStructNil is returned when the struct pointer being validated is nil.
var ErrStructNil = NewError(CodeStructNil, "the object being validated must not be nil")

// Validate runs the rules in order and short-circuits on the first failure,
// returning that error. Returns nil if all rules pass.
func Validate[T any](value T, rules ...Rule[T]) error {
	for _, rule := range rules {
		if err := rule(value); err != nil {
			return err
		}
	}
	return nil
}

// Skip is a rule that always passes, used to explicitly skip validation,
// e.g. as a placeholder when assembling rule lists conditionally.
func Skip[T any](value T) error {
	return nil
}

// FieldRule is the abstraction of a single field validation task, implemented by FieldBuilder.
// The interface method is unexported, so external packages cannot implement it
// themselves and must use the builders provided by the library.
type FieldRule interface {
	fieldRule(base reflect.Value) error
}

// FieldBuilder configures the validation task of a single field with a chain of options.
// The generic type parameter is specified explicitly by the constructor, so that
// WithField and WithRules agree on the value type at compile time.
//
//	verax.Field[string]().
//	    WithField(&user.Name).
//	    WithRules(rules.Required, rules.Length[string](2, 32)).
//	    WithLabel("Name")
type FieldBuilder[T any] struct {
	ptr   *T
	rules []Rule[T]
	// label is the field label; once set, messages are rendered as label prefix + original message
	label string
	// customErr is a custom error message; once set, any validation failure of this field returns this message
	customErr string
	// checks are parameterless custom validation functions, run in order with rules
	checks []func() error
	// fieldCmps are cross-field comparisons, run after the rule chain in order
	fieldCmps []fieldCmp[T]
}

// fieldCmp describes one cross-field comparison; code is the error code on failure.
type fieldCmp[T any] struct {
	other *T
	code  string
}

// Field creates a field validation builder, T is the field value type.
func Field[T any]() *FieldBuilder[T] {
	return &FieldBuilder[T]{}
}

// WithField binds the field pointer; the external field name is located via reflection during validation.
func (b *FieldBuilder[T]) WithField(ptr *T) *FieldBuilder[T] {
	b.ptr = ptr
	return b
}

// WithRules appends a chain of validation rules, executed in order with short-circuit.
func (b *FieldBuilder[T]) WithRules(rules ...Rule[T]) *FieldBuilder[T] {
	b.rules = append(b.rules, rules...)
	return b
}

// WithLabel sets the field label; once set, messages are rendered as label prefix + original message.
// For example, label "Name" + rule message "cannot be blank" = "Name cannot be blank";
// if not set, the original message is kept.
func (b *FieldBuilder[T]) WithLabel(label string) *FieldBuilder[T] {
	b.label = label
	return b
}

// WithErr sets a custom error message; once set, any validation failure of this field returns this message.
func (b *FieldBuilder[T]) WithErr(message string) *FieldBuilder[T] {
	b.customErr = message
	return b
}

// WithCheckFn appends a parameterless custom validation function; returning non-nil fails, run in order with rules.
func (b *FieldBuilder[T]) WithCheckFn(fn func() error) *FieldBuilder[T] {
	b.checks = append(b.checks, fn)
	return b
}

// WithFieldEq requires the value of this field to equal the value pointed to by other.
func (b *FieldBuilder[T]) WithFieldEq(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldEq)
}

// WithFieldNe requires the value of this field to differ from the value pointed to by other.
func (b *FieldBuilder[T]) WithFieldNe(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldNe)
}

// WithFieldGt requires the value of this field to be strictly greater than the value pointed to by other, ordered types only.
func (b *FieldBuilder[T]) WithFieldGt(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldGt)
}

// WithFieldGte requires the value of this field to be greater than or equal to the value pointed to by other, ordered types only.
func (b *FieldBuilder[T]) WithFieldGte(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldGte)
}

// WithFieldLt requires the value of this field to be strictly less than the value pointed to by other, ordered types only.
func (b *FieldBuilder[T]) WithFieldLt(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldLt)
}

// WithFieldLte requires the value of this field to be less than or equal to the value pointed to by other, ordered types only.
func (b *FieldBuilder[T]) WithFieldLte(other *T) *FieldBuilder[T] {
	return b.withFieldCompare(other, codes.CodeFieldLte)
}

func (b *FieldBuilder[T]) withFieldCompare(other *T, code string) *FieldBuilder[T] {
	if other != nil {
		b.fieldCmps = append(b.fieldCmps, fieldCmp[T]{other: other, code: code})
	}
	return b
}

// fieldRule implements FieldRule: runs the validation of a single field and returns the error
// (including field attribution and configuration application).
func (b *FieldBuilder[T]) fieldRule(base reflect.Value) error {
	if b.ptr == nil {
		return nil
	}
	fv := reflect.ValueOf(b.ptr).Elem()

	// A pointer-typed field whose value is nil is treated as unfilled and skipped:
	// no rules run then, including Required;
	// for required semantics bind a non-pointer field, or dereference to a value type first
	if fv.Kind() == reflect.Pointer && fv.IsNil() {
		return nil
	}

	name := locateFieldName(base, b.ptr)
	if err := b.run(fv.Interface().(T)); err != nil {
		return WithField(b.applyConfig(err), name)
	}
	if err := b.runFieldCmps(base, name); err != nil {
		return err
	}
	return nil
}

// runFieldCmps runs the cross-field comparisons in order, returning on the first failure.
func (b *FieldBuilder[T]) runFieldCmps(base reflect.Value, name string) error {
	for _, cmp := range b.fieldCmps {
		if err := b.applyFieldCmp(base, cmp, name); err != nil {
			return err
		}
	}
	return nil
}

// applyFieldCmp runs one cross-field comparison, attributing failure to this field and applying field-level config.
func (b *FieldBuilder[T]) applyFieldCmp(base reflect.Value, cmp fieldCmp[T], name string) error {
	current := reflect.ValueOf(b.ptr).Elem().Interface().(T)
	other := reflect.ValueOf(cmp.other).Elem().Interface().(T)
	if compareValues(current, other, cmp.code) {
		return nil
	}
	otherName := locateFieldName(base, cmp.other)
	err := NewMessage(cmp.code, map[string]string{"other": otherName})
	return WithField(b.applyConfig(err), name)
}

// compareValues compares two values according to the error code semantics;
// ordered comparison only supports cmp.Ordered types, passing any other type is treated as a programming error and panics.
func compareValues(a, b any, code string) bool {
	switch code {
	case codes.CodeFieldEq:
		return reflect.DeepEqual(a, b)
	case codes.CodeFieldNe:
		return !reflect.DeepEqual(a, b)
	}
	c := orderedCompare(reflect.ValueOf(a), reflect.ValueOf(b))
	switch code {
	case codes.CodeFieldGt:
		return c > 0
	case codes.CodeFieldGte:
		return c >= 0
	case codes.CodeFieldLt:
		return c < 0
	default: // CodeFieldLte
		return c <= 0
	}
}

// orderedCompare compares two ordered values and returns a negative/zero/positive number,
// with the same semantics as cmp.Compare.
func orderedCompare(a, b reflect.Value) int {
	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(a.Int(), b.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return cmp.Compare(a.Uint(), b.Uint())
	case reflect.Float32, reflect.Float64:
		return cmp.Compare(a.Float(), b.Float())
	case reflect.String:
		return cmp.Compare(a.String(), b.String())
	default:
		panic(fmt.Sprintf("verax: field comparison requires an ordered type, got %s", a.Kind()))
	}
}

// run runs the field's custom checks and rule chain in order, short-circuiting on any failure.
func (b *FieldBuilder[T]) run(value T) error {
	for _, check := range b.checks {
		if err := check(); err != nil {
			return err
		}
	}
	for _, rule := range b.rules {
		if err := rule(value); err != nil {
			return err
		}
	}
	return nil
}

// applyConfig applies field-level config to the error:
// overrides the message when a custom error is set; renders label prefix + original message when a label is set.
func (b *FieldBuilder[T]) applyConfig(err error) error {
	if b.customErr != "" {
		return &Error{Code: codeOf(err), Message: b.customErr}
	}
	if b.label != "" {
		return withLabelPrefix(err, b.label)
	}
	return err
}

// withLabelPrefix prepends the label prefix to the error message;
// non-*Error errors are wrapped by labelWrapped to preserve the original error chain, so errors.Is/As keep working.
func withLabelPrefix(err error, label string) error {
	if e, ok := errors.AsType[*Error](err); ok {
		return &Error{Code: e.Code, Message: label + e.Message, Field: e.Field}
	}
	return &labelWrapped{inner: err, prefix: label}
}

// codeOf extracts the error code; returns an empty string for non-*Error errors.
func codeOf(err error) string {
	if e, ok := errors.AsType[*Error](err); ok {
		return e.Code
	}
	return ""
}

// ValidateStruct validates each field in order and aggregates all errors.
// obj must be a non-nil pointer, otherwise ErrStructNil is returned.
// Field pointers must point to exported fields of obj, otherwise it is treated as a programming error and panics.
// Returns nil when all fields pass, or Errors when any field fails, keeping the field declaration order.
func ValidateStruct[S any](obj *S, fields ...FieldRule) error {
	if obj == nil {
		return ErrStructNil
	}
	if len(fields) == 0 {
		return nil
	}

	base := reflect.ValueOf(obj).Elem()
	var errs Errors

	for _, field := range fields {
		if err := field.fieldRule(base); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
