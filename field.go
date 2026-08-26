package verax

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// fieldMeta caches the external name of each exported field of a struct type,
// derived as: json tag first, falling back to the snake_case Go field name when missing or "-".
type fieldMeta struct {
	index int    // field index, used for address location
	name  string // external name
}

var fieldCache sync.Map // reflect.Type -> []fieldMeta

// fieldsOf returns the field metadata of a struct type, cached per type.
func fieldsOf(t reflect.Type) []fieldMeta {
	if v, ok := fieldCache.Load(t); ok {
		return v.([]fieldMeta)
	}

	infos := make([]fieldMeta, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// Unexported fields cannot be bound by external packages; zero-size fields share one address across instances and cannot be located reliably
		if !f.IsExported() || f.Type.Size() == 0 {
			continue
		}
		infos = append(infos, fieldMeta{
			index: i,
			name:  fieldNameFromTag(f.Tag.Get("json"), f.Name),
		})
	}

	v, _ := fieldCache.LoadOrStore(t, infos)
	return v.([]fieldMeta)
}

// fieldNameFromTag extracts the name from a json tag (before the first comma),
// falling back to the snake_case field name when missing or "-".
func fieldNameFromTag(tag, goName string) string {
	if i := strings.Index(tag, ","); i >= 0 {
		tag = tag[:i]
	}
	tag = strings.TrimSpace(tag)
	if len(tag) > 0 && tag != "-" {
		return tag
	}
	return Snake(goName)
}

// locateFieldName determines the field that ptr belongs to by address matching and returns its external name.
// base is a dereferenced struct value; ptr must point to an exported field of base,
// otherwise it is treated as a programming error and panics.
func locateFieldName(base reflect.Value, ptr any) string {
	target := reflect.ValueOf(ptr).Pointer()

	for _, meta := range fieldsOf(base.Type()) {
		if base.Field(meta.index).Addr().Pointer() == target {
			return meta.name
		}
	}
	// Zero-size fields share one address across instances and cannot be located reliably, so give an explicit hint
	for i := 0; i < base.NumField(); i++ {
		f := base.Type().Field(i)
		if f.Type.Size() == 0 && base.Field(i).Addr().Pointer() == target {
			panic(fmt.Sprintf("verax: field %s has zero size, its address is shared and cannot locate the field", f.Name))
		}
	}
	panic("verax: the field pointer does not belong to the struct passed to ValidateStruct")
}
