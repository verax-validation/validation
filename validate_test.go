package verax_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/verax-validation/validation"
	"github.com/verax-validation/validation/collections"
	"github.com/verax-validation/validation/is"
	"github.com/verax-validation/validation/rules"
)

var errFail = errors.New("fail")

// genericNotBlank is a generic rule function used to verify compiler argument inference
func genericNotBlank[T ~string](v T) error {
	if len(v) == 0 {
		return errFail
	}
	return nil
}

type named string

type onlyZeroSized struct {
	Empty struct{}
}

func TestValidateStructZeroSizedField(t *testing.T) {
	// a struct with only a zero-size field: the field address cannot be located reliably, so an explicit panic hint is expected
	z := &onlyZeroSized{}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "zero size") {
			t.Errorf("panic message = %v, want mention of zero size", r)
		}
	}()

	_ = verax.ValidateStruct(z,
		verax.Field[struct{}]().WithField(&z.Empty).WithRules(rules.Required[struct{}]),
	)
}

func TestValidateAllPassed(t *testing.T) {
	err := verax.Validate("abc", genericNotBlank[string])

	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestValidateShortCircuit(t *testing.T) {
	secondCalled := false
	second := func(v string) error {
		secondCalled = true
		return nil
	}

	err := verax.Validate("", genericNotBlank, second)

	if !errors.Is(err, errFail) {
		t.Errorf("Validate() = %v, want %v", err, errFail)
	}
	if secondCalled {
		t.Error("second rule should not run after first failure")
	}
}

func TestValidateMixRuleForms(t *testing.T) {
	var declared verax.Rule[string] = func(v string) error { return nil }
	literal := func(v string) error { return nil }

	err := verax.Validate("abc",
		genericNotBlank,
		declared,
		literal,
		func(v string) error { return nil },
	)

	if err != nil {
		t.Errorf("mixed rule forms failed: %v", err)
	}
}

func TestSkipAlwaysPasses(t *testing.T) {
	if err := verax.Validate("anything", verax.Skip); err != nil {
		t.Errorf("Skip should always pass, got %v", err)
	}
	if err := verax.Validate(42, verax.Skip); err != nil {
		t.Errorf("Skip should work for any type, got %v", err)
	}
}

func TestWhenConditionTrue(t *testing.T) {
	rule := verax.When[string](true, genericNotBlank)

	if err := rule(""); !errors.Is(err, errFail) {
		t.Errorf("When(true) = %v, want %v", err, errFail)
	}
}

func TestWhenConditionFalse(t *testing.T) {
	called := false
	failing := func(v string) error {
		called = true
		return errFail
	}
	rule := verax.When(false, failing)

	if err := rule(""); err != nil {
		t.Errorf("When(false) = %v, want nil", err)
	}
	if called {
		t.Error("rules should not run when condition is false")
	}
}

func TestValidateStructAggregatesFields(t *testing.T) {
	type User struct {
		Name  string
		Email string
	}
	user := &User{Name: "", Email: "not-an-email"}

	err := verax.ValidateStruct(user,
		verax.Field[string]().WithField(&user.Name).WithRules(genericNotBlank[string]),
		verax.Field[string]().WithField(&user.Email).WithRules(genericNotBlank[string], func(v string) error {
			if strings.Contains(v, "@") {
				return nil
			}
			return errors.New("must be a valid email address")
		}),
	)

	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("ValidateStruct() = %T, want verax.Errors", err)
	}
	if len(errs) != 2 {
		t.Fatalf("error count = %d, want 2", len(errs))
	}
	for _, name := range []string{"name", "email"} {
		if _, found := errs.Get(name); !found {
			t.Errorf("missing field %q in %v", name, errs)
		}
	}
}

func TestValidateStructAllPassed(t *testing.T) {
	type User struct {
		Name string
	}
	user := &User{Name: "alice"}

	err := verax.ValidateStruct(user,
		verax.Field[string]().WithField(&user.Name).WithRules(genericNotBlank[string]),
	)

	if err != nil {
		t.Errorf("ValidateStruct() = %v, want nil", err)
	}
}

func TestValidateStructNilObject(t *testing.T) {
	type User struct{}

	err := verax.ValidateStruct[User](nil)

	if !errors.Is(err, verax.ErrStructNil) {
		t.Errorf("ValidateStruct(nil) = %v, want ErrStructNil", err)
	}
}

func TestFieldNilPointerSkipped(t *testing.T) {
	type Form struct {
		Nickname *string
	}
	form := &Form{}

	err := verax.ValidateStruct(form,
		// nil pointer fields are skipped at runtime by ValidateStruct, the closure should not be called
		verax.Field[*string]().WithField(&form.Nickname).WithRules(func(v *string) error {
			return errors.New("should not run for nil field")
		}),
	)

	if err != nil {
		t.Errorf("nil field pointer should be skipped, got %v", err)
	}
}

// nestedValidatable verifies that a nested object participates in outer validation through Valid
type Address struct {
	City string
}

func (a Address) Validate() error {
	return verax.Validate(a.City, genericNotBlank)
}

type Member struct {
	Name    string
	Address Address
	Age     int
}

func TestNestedValidatable(t *testing.T) {
	member := &Member{Name: "alice"}

	err := verax.ValidateStruct(member,
		verax.Field[string]().WithField(&member.Name).WithRules(genericNotBlank[string]),
		verax.Field[Address]().WithField(&member.Address).WithRules(verax.Valid[Address]),
	)

	errs, ok := errors.AsType[verax.Errors](err)
	if !ok {
		t.Fatalf("ValidateStruct() = %T, want verax.Errors", err)
	}
	if _, found := errs.Get("address"); !found {
		t.Errorf("nested failure should be reported under address, got %v", errs)
	}
}

func TestNestedValidatablePassed(t *testing.T) {
	member := &Member{Name: "alice"}
	member.Address.City = "hangzhou"

	err := verax.ValidateStruct(member,
		verax.Field[string]().WithField(&member.Name).WithRules(genericNotBlank[string]),
		verax.Field[Address]().WithField(&member.Address).WithRules(verax.Valid[Address]),
	)

	if err != nil {
		t.Errorf("ValidateStruct() = %v, want nil", err)
	}
}

func TestNamedStringTypeInference(t *testing.T) {
	code := named("GO")

	err := verax.Validate(code, genericNotBlank[named])

	if err != nil {
		t.Errorf("named type validation failed: %v", err)
	}
}

func TestMember(t *testing.T) {
	verax.RegisterZhCN()
	defer verax.RegisterEn()

	member := Member{
		Name:    "",
		Address: Address{},
		Age:     9,
	}

	nameFiled := verax.Field[string]().WithRules(rules.Required[string], rules.Length[string](1, 5))

	err := verax.ValidateStruct(&member,
		// generic rules are inferred automatically when inlined, no explicit [string] needed
		nameFiled.WithField(&member.Name),
		verax.Field[int]().WithField(&member.Age).WithRules(rules.Between(11, 18)).WithErr("xxxxxxx"),
	)
	if err == nil {
		t.Fatal("expected failure for blank name")
	}

	got := err.Error()
	t.Logf("output: %s", got)

	if want := "name: 不能为空; age: xxxxxxx"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

}

// ---- 完整规则使用示例 ----

func TestUsageRequiredAndLength(t *testing.T) {
	// Required: 值不是其类型的零值
	if err := verax.Validate("", rules.Required); err == nil {
		t.Error("Required(empty) should fail")
	}
	if err := verax.Validate("alice", rules.Required); err != nil {
		t.Errorf("Required(alice) = %v, want nil", err)
	}

	// Optional: 零值放行, 非零值才执行后续规则
	opt := rules.Optional(rules.Length[string](2, 8))
	if err := opt(""); err != nil {
		t.Errorf("Optional(empty) = %v, want nil", err)
	}
	if err := opt("x"); err == nil {
		t.Error("Optional with out-of-range value should fail")
	}

	// RequiredIf: 条件为真时必填
	required := true
	if err := verax.Validate("", rules.RequiredIf[string](func() bool { return required })); err == nil {
		t.Error("RequiredIf(true)(empty) should fail")
	}
	required = false
	if err := verax.Validate("", rules.RequiredIf[string](func() bool { return required })); err != nil {
		t.Errorf("RequiredIf(false)(empty) = %v, want nil", err)
	}

	// Length: 字节长度落在闭区间
	if err := verax.Validate("abc", rules.Length[string](2, 3)); err != nil {
		t.Errorf("Length(2,3)(abc) = %v, want nil", err)
	}
	if err := verax.Validate("a", rules.Length[string](2, 3)); err == nil {
		t.Error("Length(2,3)(a) should fail")
	}

	// Len: 字节长度精确等于 n
	if err := verax.Validate("abc", rules.Len[string](3)); err != nil {
		t.Errorf("Len(3)(abc) = %v, want nil", err)
	}
	if err := verax.Validate("ab", rules.Len[string](3)); err == nil {
		t.Error("Len(3)(ab) should fail")
	}

	// NotNil: 指针非 nil
	var p *string
	if err := verax.Validate(p, rules.NotNil); err == nil {
		t.Error("NotNil(nil) should fail")
	}
	name := "x"
	if err := verax.Validate(&name, rules.NotNil); err != nil {
		t.Errorf("NotNil(&x) = %v, want nil", err)
	}
}

func TestUsageNumericAndEnum(t *testing.T) {
	// Min / Max / Between: 闭区间
	if err := verax.Validate(18, rules.Min(18)); err != nil {
		t.Errorf("Min(18)(18) = %v, want nil", err)
	}
	if err := verax.Validate(17, rules.Min(18)); err == nil {
		t.Error("Min(18)(17) should fail")
	}
	if err := verax.Validate(100, rules.Max(100)); err != nil {
		t.Errorf("Max(100)(100) = %v, want nil", err)
	}
	if err := verax.Validate(101, rules.Max(100)); err == nil {
		t.Error("Max(100)(101) should fail")
	}
	if err := verax.Validate(5, rules.Between(1, 10)); err != nil {
		t.Errorf("Between(1,10)(5) = %v, want nil", err)
	}
	if err := verax.Validate(11, rules.Between(1, 10)); err == nil {
		t.Error("Between(1,10)(11) should fail")
	}

	// Gt / Lt: 严格区间
	if err := verax.Validate(19, rules.Gt(18)); err != nil {
		t.Errorf("Gt(18)(19) = %v, want nil", err)
	}
	if err := verax.Validate(18, rules.Gt(18)); err == nil {
		t.Error("Gt(18)(18) should fail")
	}
	if err := verax.Validate(59, rules.Lt(60)); err != nil {
		t.Errorf("Lt(60)(59) = %v, want nil", err)
	}
	if err := verax.Validate(60, rules.Lt(60)); err == nil {
		t.Error("Lt(60)(60) should fail")
	}

	// Eq / Ne: 等值
	if err := verax.Validate("admin", rules.Eq("admin")); err != nil {
		t.Errorf("Eq(admin)(admin) = %v, want nil", err)
	}
	if err := verax.Validate("user", rules.Eq("admin")); err == nil {
		t.Error("Eq(admin)(user) should fail")
	}
	if err := verax.Validate("user", rules.Ne("admin")); err != nil {
		t.Errorf("Ne(admin)(user) = %v, want nil", err)
	}
	if err := verax.Validate("admin", rules.Ne("admin")); err == nil {
		t.Error("Ne(admin)(admin) should fail")
	}

	// MultipleOf: 整数倍
	if err := verax.Validate(100, rules.MultipleOf(10)); err != nil {
		t.Errorf("MultipleOf(10)(100) = %v, want nil", err)
	}
	if err := verax.Validate(5, rules.MultipleOf(10)); err == nil {
		t.Error("MultipleOf(10)(5) should fail")
	}

	// In / NotIn: 枚举
	if err := verax.Validate("guest", rules.In("guest", "admin")); err != nil {
		t.Errorf("In(guest,admin)(guest) = %v, want nil", err)
	}
	if err := verax.Validate("hacker", rules.In("guest", "admin")); err == nil {
		t.Error("In(guest,admin)(hacker) should fail")
	}
	if err := verax.Validate("guest", rules.NotIn("banned")); err != nil {
		t.Errorf("NotIn(banned)(guest) = %v, want nil", err)
	}
	if err := verax.Validate("banned", rules.NotIn("banned")); err == nil {
		t.Error("NotIn(banned)(banned) should fail")
	}
}

func TestUsageStringRules(t *testing.T) {
	// Match: 正则匹配
	if err := verax.Validate("ABC123", rules.Match(`^[A-Z0-9]+$`)); err != nil {
		t.Errorf("Match(upper)(ABC123) = %v, want nil", err)
	}
	if err := verax.Validate("abc", rules.Match(`^[A-Z0-9]+$`)); err == nil {
		t.Error("Match(upper)(abc) should fail")
	}

	// Date: 符合指定时间布局
	if err := verax.Validate("2026-08-26", rules.Date("2006-01-02")); err != nil {
		t.Errorf("Date(2006-01-02)(valid) = %v, want nil", err)
	}
	if err := verax.Validate("2026/08/26", rules.Date("2006-01-02")); err == nil {
		t.Error("Date(2006-01-02)(slash) should fail")
	}

	// Contains / ContainsAny / Excludes
	if err := verax.Validate("alice@example.com", rules.Contains("@")); err != nil {
		t.Errorf("Contains(@) = %v, want nil", err)
	}
	if err := verax.Validate("alice@example.com", rules.ContainsAny("!@#")); err != nil {
		t.Errorf("ContainsAny(!@#) = %v, want nil", err)
	}
	if err := verax.Validate("alice@example.com", rules.Excludes(" ")); err != nil {
		t.Errorf("Excludes(space) = %v, want nil", err)
	}

	// StartWith / EndWith
	if err := verax.Validate("https://example.com", rules.StartWith("https://")); err != nil {
		t.Errorf("StartWith(https://) = %v, want nil", err)
	}
	if err := verax.Validate("image.png", rules.EndWith(".png")); err != nil {
		t.Errorf("EndWith(.png) = %v, want nil", err)
	}

	// TrimSpace: 容忍输入首尾空白
	trimmed := rules.TrimSpace(rules.Required)
	if err := trimmed("  alice  "); err != nil {
		t.Errorf("TrimSpace(Required)(spaced) = %v, want nil", err)
	}
	if err := trimmed("   "); err == nil {
		t.Error("TrimSpace(Required)(blank) should fail")
	}
}

func TestUsageFormatRules(t *testing.T) {
	// is 包常用格式断言: 合法值逐一通过
	valid := []struct {
		name string
		rule func(string) error
		val  string
	}{
		{"EmailFormat", is.EmailFormat, "alice@example.com"},
		{"URL", is.URL, "https://example.com"},
		{"IP", is.IP, "192.168.1.1"},
		{"IPv4", is.IPv4, "192.168.1.1"},
		{"UUIDv4", is.UUIDv4, "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
		{"JSON", is.JSON, `{"a":1}`},
		{"Base64", is.Base64, "aGVsbG8="},
		{"CreditCard", is.CreditCard, "4111111111111111"},
		{"MobileCN", is.MobileCN, "13812345678"},
		{"IDCardCN", is.IDCardCN, "11010519491231002X"},
		{"PostalCodeCN", is.PostalCodeCN, "310000"},
		{"CountryCode", is.CountryCode, "CN"},
		{"CurrencyCode", is.CurrencyCode, "CNY"},
		{"LanguageCode", is.LanguageCode, "zh"},
		{"CIDR", is.CIDR, "192.168.0.0/24"},
		{"JWT", is.JWT, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
		{"Boolean", is.Boolean, "true"},
		{"TimeZone", is.TimeZone, "Asia/Shanghai"},
	}
	for _, tc := range valid {
		if err := tc.rule(tc.val); err != nil {
			t.Errorf("%s(%q) = %v, want nil", tc.name, tc.val, err)
		}
	}

	// 非法值示例
	invalid := []struct {
		name string
		rule func(string) error
		val  string
	}{
		{"EmailFormat", is.EmailFormat, "plain"},
		{"IP", is.IP, "999.1.1.1"},
		{"MobileCN", is.MobileCN, "12345678901"},
		{"CountryCode", is.CountryCode, "XX"},
		{"CIDR", is.CIDR, "192.168.0.0"},
	}
	for _, tc := range invalid {
		if err := tc.rule(tc.val); err == nil {
			t.Errorf("%s(%q) should fail", tc.name, tc.val)
		}
	}
}

func TestUsageCollectionRules(t *testing.T) {
	// SliceLen: 切片长度闭区间
	if err := verax.Validate([]string{"go", "web"}, collections.SliceLen[string](1, 5)); err != nil {
		t.Errorf("SliceLen(1,5) = %v, want nil", err)
	}

	// Unique: 元素互不重复
	if err := verax.Validate([]string{"a", "b", "c"}, collections.Unique[string]()); err != nil {
		t.Errorf("Unique(distinct) = %v, want nil", err)
	}
	if err := verax.Validate([]string{"a", "b", "a"}, collections.Unique[string]()); err == nil {
		t.Error("Unique(duplicated) should fail")
	}

	// Slice: 逐元素校验
	sliceRule := collections.Slice[string](rules.Required)
	if err := sliceRule([]string{"go", "web"}); err != nil {
		t.Errorf("Slice(Required) = %v, want nil", err)
	}
	if err := sliceRule([]string{"go", ""}); err == nil {
		t.Error("Slice(Required) with blank element should fail")
	}

	// Map: 值校验
	if err := verax.Validate(map[string]string{"en": "go"}, collections.Map[string, string](rules.Required)); err != nil {
		t.Errorf("Map(Required) = %v, want nil", err)
	}

	// MapLen: 键值对数量
	if err := verax.Validate(map[string]int{"a": 1}, collections.MapLen[string, int](1, 1)); err != nil {
		t.Errorf("MapLen(1,1) = %v, want nil", err)
	}

	// Each: iter.Seq 序列逐元素校验
	if err := verax.Validate(slices.Values([]string{"go", "web"}), collections.Each[string](rules.Required)); err != nil {
		t.Errorf("Each(Required) = %v, want nil", err)
	}
}

func TestUsageCrossFieldRules(t *testing.T) {
	type Form struct {
		Password string
		Confirm  string
		Start    int
		End      int
	}

	// 全部通过: 确认密码一致, 结束时间不早于开始时间
	form := &Form{Password: "abc", Confirm: "abc", Start: 1, End: 2}
	if err := verax.ValidateStruct(form,
		verax.Field[string]().WithField(&form.Confirm).WithFieldEq(&form.Password),
		verax.Field[int]().WithField(&form.End).WithFieldGte(&form.Start),
	); err != nil {
		t.Errorf("valid cross-field form = %v, want nil", err)
	}

	// 失败: 确认密码与密码不一致
	form = &Form{Password: "abc", Confirm: "abd", Start: 1, End: 2}
	if err := verax.ValidateStruct(form,
		verax.Field[string]().WithField(&form.Confirm).WithFieldEq(&form.Password),
	); err == nil {
		t.Error("mismatched confirm should fail")
	}

	// 失败: 结束时间早于开始时间
	form = &Form{Password: "abc", Confirm: "abc", Start: 5, End: 3}
	if err := verax.ValidateStruct(form,
		verax.Field[int]().WithField(&form.End).WithFieldGte(&form.Start),
	); err == nil {
		t.Error("end < start should fail")
	}
}
