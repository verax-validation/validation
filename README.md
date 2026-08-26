# verax

[![CI](https://github.com/verax-validation/validation/actions/workflows/ci.yml/badge.svg)](https://github.com/verax-validation/validation/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/verax-validation/validation.svg)](https://pkg.go.dev/github.com/verax-validation/validation)
[![Go Version](https://img.shields.io/github/go-mod/go-version/verax-validation/validation)](https://go.dev/dl/)
[![Release](https://img.shields.io/github/v/release/verax-validation/validation)](https://github.com/verax-validation/validation/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/verax-validation/validation)](https://goreportcard.com/report/github.com/verax-validation/validation)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**English** | [简体中文](README.zh-CN.md)

A data validation library for Go built on Go 1.26 generics. Rules are functions, and type safety is guaranteed by the compiler.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Built-in Rules](#built-in-rules)
- [Caveats](#caveats)
- [Error Handling](#error-handling)
- [Internationalization](#internationalization)
- [Comparison with Other Libraries](#comparison-with-other-libraries)
- [Versioning and Stability](#versioning-and-stability)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Fully generic API**: `Rule[T]` binds rules to value types at compile time; passing a mismatched rule fails to compile
- **Chained field configuration**: `Field[T]()` binds pointers, appends rules, and sets labels or custom errors via an option pattern
- **Cross-field comparison**: `WithFieldEq` / `WithFieldGte` and friends express "confirm password must match" or "end time ≥ start time" in one line
- **Conditional required**: `rules.RequiredIf` declaratively expresses constraints like "card number is required once an online payment method is selected"
- **High performance**: per-element slice validation is an order of magnitude faster than reflection-based implementations, see [benchmark report](docs/benchmark.md)
- **Explicit optional semantics**: rules always run by default; optional fields are declared explicitly with `Optional`, no hidden behavior
- **Process-level internationalization**: six built-in languages, templates keep placeholders, rule parameters are interpolated at render time
- **Zero third-party dependencies**: built entirely on the standard library
- **Structured errors**: error code + message template interpolation + field context, extractable via `errors.AsType`

## Installation

```
go get github.com/verax-validation/validation
```

Requires Go 1.26 or above.

## Getting Started

### Validating Simple Values

```go
import (
    "github.com/verax-validation/validation"
    "github.com/verax-validation/validation/rules"
)

err := verax.Validate("alice@example.com",
    rules.Required,
    rules.Length[string](5, 100),
)
// bare generic rules infer their type at the call site; e.g. rules.Required needs no explicit rules.Required[string]
```

### Validating a Struct

```go
type User struct {
    Name  string
    Email string
    Age   int
    Bio   string
}

func (u *User) Validate() error {
    return verax.ValidateStruct(u,
        verax.Field[string]().
            WithField(&u.Name).
            WithRules(rules.Required, rules.Length[string](2, 64)),
        verax.Field[string]().
            WithField(&u.Email).
            WithRules(rules.Required, is.EmailFormat).
            WithLabel("Email"),
        // optional field: empty value skips the following rules
        verax.Field[string]().
            WithField(&u.Bio).
            WithRules(rules.Optional(rules.Length[string](0, 200))),
    )
}
```

Fields are configured with a chained option pattern: `WithField` binds a pointer, `WithRules` appends rules,
`WithLabel` sets a label (error prefix, e.g. "Email cannot be blank"), `WithErr` sets a custom error,
`WithCheckFn` adds a parameterless custom check. A mismatch between the element type and the rule value type fails to compile.

### Cross-field Comparison

```go
type RegisterForm struct {
    Password string
    Confirm  string
    StartAt  int
    EndAt    int
}

err := verax.ValidateStruct(&form,
    // confirm password must match the password
    verax.Field[string]().WithField(&form.Confirm).WithFieldEq(&form.Password),
    // end time must not be earlier than start time
    verax.Field[int]().WithField(&form.EndAt).WithFieldGte(&form.StartAt),
)
// the failure is attributed to the current field and references the other field name: "confirm: must equal password"
```

`WithFieldEq` / `WithFieldNe` / `WithFieldGt` / `WithFieldGte` / `WithFieldLt` / `WithFieldLte`
cover equality and ordered comparison, and also compose with `WithLabel` / `WithErr`.

### Conditional Required

```go
type Order struct {
    PaymentMethod string
    CardNumber    string
}

err := verax.ValidateStruct(&order,
    verax.Field[string]().
        WithField(&order.CardNumber).
        WithRules(
            // the card number is required when an online payment method is selected
            rules.RequiredIf[string](func() bool { return order.PaymentMethod == "card" }),
            is.CreditCard,
        ),
)
```

### Nested Objects

```go
type Address struct {
    City string
}

func (a Address) Validate() error {
    return verax.Validate(a.City, rules.Required)
}

err := verax.ValidateStruct(user,
    verax.Field[Address]().WithField(&user.Address).WithRules(verax.Valid),
)
```

### Validating Collections

```go
import "github.com/verax-validation/validation/collections"

err := verax.Validate(tags,
    collections.SliceLen[string](1, 10),            // size constraint
    collections.Unique[string](),                   // elements must be distinct
    collections.Slice[string](                      // per-element validation
        rules.Required,
        rules.Length[string](1, 20),
    ),
)
// element errors are aggregated by position: Errors{"1": cannot be blank}
```

`collections` provides `Slice` / `Each`(iter.Seq) / `Map` for per-element validation, plus
`SliceLen` / `MapLen` for size and `Unique` for distinctness.

### Conditional Rules

```go
// isAdmin is a condition already determined on the business side, evaluated when the rule is constructed
rule := verax.When(isAdmin, rules.Min(-1))
```

### Complete Example

A fully runnable example covering struct validation, cross-field comparison,
collection validation, Chinese localization and error handling is in [examples/basic](examples/basic):

```
go run ./examples/basic
```

## Built-in Rules

### rules package — General Value Rules

Work with any type, used through `verax.Validate` or `Field.WithRules`.

**Required and Optional**

| Rule | Description |
|------|-------------|
| `Required` | The value is not the zero value of its type (empty string / 0 / false / nil pointer all count as empty) |
| `RequiredIf(cond)` | The value is required when the condition is true |
| `Optional(rules...)` | Zero values pass directly, non-zero values run the given rules |

**Length** (strings and byte slices)

| Rule | Description |
|------|-------------|
| `Length(min, max)` | Byte length falls within the closed interval `[min, max]` |
| `Len(n)` | Byte length is exactly n |

**Numeric and Range** (ordered types: integers / floats / strings / time, etc.)

| Rule | Description |
|------|-------------|
| `Min(min)` | Greater than or equal to min |
| `Max(max)` | Less than or equal to max |
| `Between(lo, hi)` | Falls within the closed interval `[lo, hi]` |
| `Gt(min)` | Strictly greater than min (complements the closed-interval Min) |
| `Lt(max)` | Strictly less than max (complements the closed-interval Max) |
| `MultipleOf(step)` | An integer multiple of step |

**Equality and Enum**

| Rule | Description |
|------|-------------|
| `Eq(target)` | Equal to the target value |
| `Ne(target)` | Not equal to the target value |
| `In(values...)` | Matches an allowed value list (hashed into a set at construction, O(1) lookup) |
| `NotIn(values...)` | Does not match a forbidden value list |

**String**

| Rule | Description |
|------|-------------|
| `Match(pattern)` | Matches a regular expression (compiled at construction; an invalid expression panics immediately) |
| `Date(layout)` | Matches a time.Parse layout, e.g. "2006-01-02" |
| `Contains(substr)` | Contains the given substring |
| `ContainsAny(chars)` | Contains any of the given characters (e.g. password complexity) |
| `Excludes(substr)` | Does not contain the given substring |
| `StartWith(prefix)` | Starts with the given prefix |
| `EndWith(suffix)` | Ends with the given suffix |
| `TrimSpace(inner)` | Trims leading/trailing whitespace before applying inner, tolerating surrounding spaces |

**Pointer**

| Rule | Description |
|------|-------------|
| `NotNil` | The pointer is not nil (for validating pointers directly with Validate, or validating pointer elements in collections) |

### is package — Format Assertions (71)

All operate on strings in strict mode (an empty string is always invalid); wrap optional fields with `rules.Optional`.

**Email**

| Rule | Description |
|------|-------------|
| `Email` | Well-formed and the domain has MX records (requires a DNS lookup) |
| `EmailFormat` | Well-formed only, no DNS lookup |

**Character Sets**

| Rule | Description |
|------|-------------|
| `Alpha` / `Alphanumeric` | English letters only / English letters and digits only |
| `Digit` | ASCII digits only |
| `UTFLetter` / `UTFDigit` / `UTFNumeric` / `UTFLetterNumeric` | unicode letters / decimal digits / number characters / letters or numbers |
| `LowerCase` / `UpperCase` | All lowercase / all uppercase |
| `ASCII` / `PrintableASCII` | ASCII only / printable ASCII only |
| `Hexadecimal` | Hexadecimal digit sequence |

**Numeric and Boolean**

| Rule | Description |
|------|-------------|
| `Int` | Integer literal (allowing a sign) |
| `Float` | Decimal floating-point literal (allowing scientific notation) |
| `Boolean` | true / false / t / f / 1 / 0, case-insensitive |

**Color**

| Rule | Description |
|------|-------------|
| `HexColor` | #RGB or #RRGGBB |
| `RGBColor` | rgb(r, g, b), components 0-255 |
| `RGBA` | rgba(r, g, b, a), a is a 0-1 decimal or 0-100% |
| `HSL` / `HSLA` | hsl / hsla, h 0-360, s/l 0-100 |

**Identifiers**

| Rule | Description |
|------|-------------|
| `UUID` / `UUIDv3` / `UUIDv4` / `UUIDv5` / `UUIDv7` | Generic / version-specific UUID |
| `ULID` | 26-character Crockford Base32 |
| `MongoID` | 24-hex-digit MongoDB ObjectID |

**Text**

| Rule | Description |
|------|-------------|
| `JSON` | A syntactically complete JSON document |
| `Base64` | Standard base64 (RFC 4648, with padding) |
| `Base64URL` | URL-safe base64 (without padding) |
| `JWT` | Three-segment JWT (each segment decodable as base64url) |
| `DataURI` | RFC 2397 data URI (including base64 payload decode check) |

**Network**

| Rule | Description |
|------|-------------|
| `URL` | Absolute URL (scheme + host) |
| `RequestURL` / `RequestURI` | Request URL / request URI |
| `IP` / `IPv4` / `IPv6` | IP addresses |
| `MAC` | MAC address (colon / hyphen / dot separators) |
| `Host` | IP address or DNS name |
| `Port` | Port 1-65535 |
| `DialString` | host:port dial string |
| `E164` | E.164 international phone number |
| `CIDR` / `CIDRv4` / `CIDRv6` | CIDR notation |
| `Subdomain` | Single subdomain label (1-63 characters) |
| `Domain` | Full domain (at least two labels) |
| `DNSName` | DNS name (trailing root dot allowed) |

**Geography**

| Rule | Description |
|------|-------------|
| `Latitude` | Latitude [-90, 90] |
| `Longitude` | Longitude [-180, 180] |

**Cards / Book Numbers**

| Rule | Description |
|------|-------------|
| `CreditCard` | A card number passing the Luhn check |
| `ISBN10` / `ISBN13` / `ISBN` | ISBN-10 / ISBN-13 / either |

**Business Formats**

| Rule | Description |
|------|-------------|
| `SSN` | US social security number |
| `Semver` | Semantic versioning 2.0.0 |
| `Origin` | CORS origin (http(s) + host + optional port) |
| `TimeZone` | A system-recognizable timezone name, e.g. Asia/Shanghai |

**Character Width**

| Rule | Description |
|------|-------------|
| `Multibyte` | Contains at least one multibyte character |
| `FullWidth` / `HalfWidth` / `VariableWidth` | Contains full-width / half-width / both |

**China Localization**

| Rule | Description |
|------|-------------|
| `MobileCN` | Mainland China mobile number (11 digits, starts with 1, second digit 3-9) |
| `IDCardCN` | 18-digit ID card number (validates birth date and GB 11643 check digit) |
| `PostalCodeCN` | 6-digit postal code |

**International Standard Codes**

| Rule | Description |
|------|-------------|
| `CountryCode` | ISO 3166-1 alpha-2 country code (complete set) |
| `CurrencyCode` | ISO 4217 currency code (complete set) |
| `LanguageCode` | ISO 639-1 language code (complete set) |

### collections package — Collection Validation

| Rule | Description |
|------|-------------|
| `Slice(rules...)` | Per-element slice validation, failures aggregated by index |
| `Each(rules...)` | Per-element validation of an iter.Seq sequence |
| `Map(rules...)` | Per-value map validation, failures aggregated by key |
| `SliceLen(min, max)` | Slice length within a closed interval |
| `MapLen(min, max)` | Map entry count within a closed interval |
| `Unique()` | Elements must be mutually distinct |

## Caveats

The following behaviors are by design; please note them before use to avoid pitfalls:

- **nil pointer fields are skipped**: when a pointer field bound to `FieldBuilder` is nil,
  no rules run, **including `Required`**. For required semantics bind a non-pointer field, or dereference to a value type first.
- **zero-size fields cannot be bound**: zero-size fields such as `struct{}` share one address across all instances,
  so they cannot be located reliably; binding one panics. Do not validate zero-size fields.
- **`is.Email` depends on DNS**: it queries the domain's MX records; use `is.EmailFormat` offline or in unit tests.
- **`rules.MultipleOf` computes via float64**: very large integers (above 2^53) may lose precision.
- **cross-field ordered comparison supports ordered types only**: `WithFieldGt/Gte/Lt/Lte` panic on unordered types (e.g. structs).
- **field name derivation**: `ValidateStruct` derives the external field name from the json tag, falling back to the snake_case Go field name when missing or `-`.
- **`rules.NotNil` does not apply to nil pointer fields**: as with the first point, a nil pointer field is already skipped;
  use `NotNil` to validate pointers directly with `Validate`, or to validate pointer elements in collections.

## Error Handling

Error codes are provided as constants (e.g. `rules.CodeLength`), so typos are exposed at compile time.
On failure, rules render template parameters through `NewMessage` in the currently active language, producing an `*Error` carrying the code:

```go
// inside a rule (Length as an example): render per the current language on failure
return verax.NewMessage(rules.CodeLength, map[string]string{"min": "5", "max": "100"})

// application custom error: construct the full message directly
err := verax.NewError("app.token.expired", "the token has expired")
```

```go
// extract structured info (Go 1.26 errors.AsType)
if e, ok := errors.AsType[*verax.Error](err); ok {
    fmt.Println(e.Code)               // validation.length
    fmt.Println(e.Message)            // the length must be between 5 and 100
    fmt.Println(e.Field)              // name (filled when validated through ValidateStruct)
}
```

Render order: active language table → built-in English table → the error code itself.
Aggregated errors render as `"name: cannot be blank; age: must be no less than 18"`.

## Internationalization

Six built-in languages: en / zh-CN / zh-TW / ja / fr / de, injected with one line:

```go
verax.RegisterZhCN()  // register Simplified Chinese and make it active immediately
// subsequent validation failures are automatically in Chinese: "name: 不能为空"
```

- Built-in languages: `RegisterEn` / `RegisterZhCN` / `RegisterZhTW` / `RegisterJa` / `RegisterFr` / `RegisterDe`
- Registration takes effect immediately: `RegisterLocale(locale, table)` directly replaces the active table
- Templates keep `{{.placeholder}}` (e.g. `"不能小于 {{.min}}"`), interpolated by rule parameters at render time,
  so every language receives the exact numeric values
- Custom languages can be added via `verax.RegisterLocale(locale, messages)`;
  use `verax.Codes()` to get the full error-code list and verify complete coverage

## Comparison with Other Libraries

| Dimension | verax | go-playground/validator | ozzo-validation |
|-----------|-------|------------------------|-----------------|
| Declaration | type-safe chained API | struct tags | chained API |
| Type safety | rules bound to value types at compile time | none (string tags parsed by reflection) | limited generics |
| Reflection overhead | field location only | full reflection per validation | low |
| Optional semantics | explicit `Optional` | `omitempty` tag | `NilOrNotEmpty` |
| Cross-field comparison | built-in `WithFieldEq/Gte` etc. | `eqfield` etc. tags | hand-written |
| Conditional required | `RequiredIf` | `required_if` etc. tags | hand-written |
| Collection validation | per-element + unique + size | third-party or hand-written | `Each` |
| Internationalization | six built-in languages | third-party library | third-party |
| Dependencies | zero third-party | zero | zero |

**Why not struct tags?** The tag approach encodes validation rules into strings resolved by reflection at runtime,
so typos only surface at runtime, and it cannot express generic rules or complex conditions.
verax makes every rule a `func(T) error`, with types guaranteed by the compiler, so rules compose, reuse, and unit-test freely.

## Versioning and Stability

- Follows [Semantic Versioning](https://semver.org/) SemVer: `X.Y.Z` corresponds to incompatible / feature / fix.
- **Error codes are an external contract**: the values in `internal/codes` must not change after release; new rules only add.
- The `0.x` phase allows incompatible changes (then incompatibility only bumps the minor version).
- Changes are recorded in [CHANGELOG.md](CHANGELOG.md).

## Contributing

Issues and PRs are welcome. The full flow for adding a new rule:

1. Define an error code constant in `internal/codes/codes.go`;
2. Add translations to all 6 language tables in `internal/messages/` (`TestBuiltinLocalesCompleteness` verifies full coverage);
3. Re-export the error code in the corresponding package's (`rules` / `is` / `collections`) `codes.go`;
4. Implement the rule and add tests (both passing and failing paths).

Local verification:

```
go test ./...
go vet ./...
```

Commit messages follow Conventional Commits (e.g. `feat(rules): add ContainsAny rule`).

## License

[MIT](LICENSE)
