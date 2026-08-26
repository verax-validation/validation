package verax

import (
	"bytes"
	"errors"
	"maps"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	msg "github.com/verax-validation/validation/internal/messages"
)

// Error represents a single validation error, the fully rendered result of a failed validation.
// Code is used for programmatic checks and multi-language translation, Message is the rendered
// message text, and Field is the external field name the error belongs to
// (auto-filled when validated through ValidateStruct).
type Error struct {
	Code    string
	Message string
	Field   string
}

// Error returns the message text.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a validation error and registers its error code.
// message is the final message text provided by the caller directly;
// rules inside the library should use NewMessage to render per the current language instead of this function.
func NewError(code, message string) *Error {
	registerCode(code)
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewMessage renders the message in the currently active language and returns a new *Error.
// The template is taken from the active language table (falling back to the built-in English table),
// interpolating params with Go template syntax {{.field}};
// rules call this constructor directly on failure, for example:
//
//	verax.NewMessage(codes.CodeLength, map[string]string{"min": "5", "max": "100"})
func NewMessage(code string, params map[string]string) *Error {
	registerCode(code)
	return &Error{
		Code:    code,
		Message: renderMessage(messageTemplate(code), params),
	}
}

// messageTemplate returns the message template for the error code:
// the current language table first, falling back to the built-in English table,
// then falling back to the error code itself for diagnosis.
func messageTemplate(code string) string {
	if tmpl := currentTable()[code]; len(tmpl) > 0 {
		return tmpl
	}
	if tmpl := msg.En[code]; len(tmpl) > 0 {
		return tmpl
	}
	return code
}

// renderMessage renders the message with Go template syntax, with placeholders like {{.min}}.
func renderMessage(tmplText string, params map[string]string) string {
	if !strings.Contains(tmplText, "{{") {
		return tmplText
	}
	tmpl, err := template.New("message").Parse(tmplText)
	if err != nil {
		return tmplText
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return tmplText
	}
	return buf.String()
}

var (
	codesMu         sync.Mutex
	registeredCodes = make(map[string]struct{})
)

func registerCode(code string) {
	codesMu.Lock()
	defer codesMu.Unlock()
	registeredCodes[code] = struct{}{}
}

// Codes returns all registered error codes in lexicographic order.
// Used for language table completeness self-checks: when any language table misses an error code,
// the corresponding error silently falls back to the English default, so coverage is recommended to be verified in CI.
func Codes() []string {
	codesMu.Lock()
	defer codesMu.Unlock()
	list := make([]string, 0, len(registeredCodes))
	for code := range registeredCodes {
		list = append(list, code)
	}
	sort.Strings(list)
	return list
}

// DefaultMessages returns a copy of the built-in English template table.
// Language packages maintaining an English table should keep it consistent with this return value entry by entry.
func DefaultMessages() map[string]string {
	snapshot := make(map[string]string, len(msg.En))
	maps.Copy(snapshot, msg.En)
	return snapshot
}

// WithField attaches field-name context to an error and returns a new error object.
// Shared sentinel errors (such as rules.ErrRequired) are written into immutable copies, so they do not pollute each other;
// non-*Error errors are lightly wrapped to carry field attribution, with the original error chain preserved via Unwrap,
// so errors.Is/As keep working.
// Called automatically when aggregating in ValidateStruct and collection rules; also usable when assembling custom aggregate errors.
func WithField(err error, name string) error {
	if e, ok := errors.AsType[*Error](err); ok {
		return &Error{Code: e.Code, Message: e.Message, Field: name}
	}
	return &fieldWrapped{inner: err, field: name}
}

// fieldCarrier indicates an error carrying field attribution.
type fieldCarrier interface {
	fieldName() string
}

func (e *Error) fieldName() string { return e.Field }

// fieldWrapped is a light wrapper attaching field attribution to a non-*Error error.
type fieldWrapped struct {
	inner error
	field string
}

func (w *fieldWrapped) Error() string { return w.inner.Error() }
func (w *fieldWrapped) Unwrap() error { return w.inner }
func (w *fieldWrapped) fieldName() string {
	return w.field
}

// labelWrapped is a light wrapper attaching a message prefix to a bare error, preserving the original chain for errors.Is/As.
type labelWrapped struct {
	inner  error
	prefix string
}

func (w *labelWrapped) Error() string { return w.prefix + w.inner.Error() }
func (w *labelWrapped) Unwrap() error { return w.inner }

// MessageMap defines the mapping from error codes to message templates; templates use Go template syntax, e.g. "length must be between {{.min}} and {{.max}}".
type MessageMap map[string]string

var (
	localeMu sync.RWMutex
	// currentMessages holds the message table of the currently active language, defaulting to English, only one table at a time.
	currentMessages MessageMap
)

func init() {
	currentMessages = msg.En
}

// currentTable returns the message table of the currently active language.
func currentTable() MessageMap {
	localeMu.RLock()
	defer localeMu.RUnlock()
	return currentMessages
}

// RegisterLocale registers a message table as the currently active language, replacing existing data directly.
// locale identifies the language (e.g. zh-CN) and does not participate in table indexing;
// the built-in English table is always available as a fallback through the msg.En constant.
func RegisterLocale(locale string, table MessageMap) {
	if len(table) == 0 {
		return
	}
	localeMu.Lock()
	defer localeMu.Unlock()
	currentMessages = table
}

// Errors aggregates multiple validation errors in validation order, as the return type of ValidateStruct and collection rules.
// The slice keeps validation order; field attribution is carried by each *Error's Field attribute,
// and when rendered, entries with a Field output "field: message", entries without output "index: message".
type Errors []error

// Error returns the aggregated error message, formatted as "field1: message1; field2: message2", keeping validation order.
func (errs Errors) Error() string {
	var sb strings.Builder
	for i, err := range errs {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(entryName(err, i))
		sb.WriteString(": ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// Get returns the first error matching the given field name (or the index of an unnamed entry); ok is false when not found.
func (errs Errors) Get(field string) (error, bool) {
	for i, err := range errs {
		if entryName(err, i) == field {
			return err, true
		}
	}
	return nil, false
}

// entryName returns the locating name of an entry: the field name when field attribution exists, otherwise the validation-order index.
func entryName(err error, index int) string {
	if c, ok := err.(fieldCarrier); ok {
		if name := c.fieldName(); len(name) > 0 {
			return name
		}
	}
	return strconv.Itoa(index)
}
