package reflects

import (
	"reflect"
	"strconv"
	"strings"
)

// RepresentType returns the representation of type
//  eg: str -> float32
// ptr type would be represented by *type
// slice type would be []type
// array type would be [size]type
// dict type would be map[typKey]typVal
// named type would be represented as <pkgPath>*<typeName>
// chan type would be chan[type]
// function type would be func([params])([returns])
// unnamed struct would be struct { fieldName fieldType }
// interface would be interface { [signatures] }
func RepresentType(t reflect.Type) string {
	// todo: implement naming for channel,function,unnamed type,interface
	if t.PkgPath() != "" {
		return t.PkgPath() + "." + t.Name()
	}

	switch {
	case t.Kind() == reflect.Ptr:
		return "*" + RepresentType(t.Elem())
	case t.Kind() == reflect.Map:
		return "map[" + RepresentType(t.Key()) + "]" + RepresentType(t.Elem())
	case t.Kind() == reflect.Slice:
		return "[]" + RepresentType(t.Elem())
	case t.Kind() == reflect.Array:
		return "[" + strconv.FormatInt(int64(t.Len()), 10) + "]" + RepresentType(t.Elem())
	case t.Kind() == reflect.Struct:
		builder := strings.Builder{}
		builder.WriteString("struct { ")
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			builder.WriteString(f.Name)
			builder.WriteString(" ")
			builder.WriteString(RepresentType(f.Type))
			builder.WriteString(" ")
			if len(f.Tag) != 0 {
				builder.WriteString("\"")
				builder.WriteString(strings.ReplaceAll(string(f.Tag), "\"", "\\\""))
				builder.WriteString("\"")
			}
			if i != t.NumField()-1 {
				builder.WriteString("; ")
			}
		}
		builder.WriteString("}")
		return builder.String()
	default:
		return t.String()
	}
}

func IsEmptyInterface(t reflect.Type) bool {
	return RepresentType(t) == "interface {}"
}
