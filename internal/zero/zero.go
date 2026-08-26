// Package zero provides zero-value detection for values of any type.
package zero

import "reflect"

// IsZero reports whether value is the zero value of its type.
// Common scalar types take the fast path, other types fall back to reflection;
// pointers are dereferenced level by level, a pointer to a zero value is considered zero, and a nil pointer is also considered zero.
func IsZero(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return len(v) == 0
	case bool:
		return !v
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	default:
		rv := reflect.ValueOf(value)
		for rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return true
			}
			rv = rv.Elem()
		}
		return rv.IsZero()
	}
}
