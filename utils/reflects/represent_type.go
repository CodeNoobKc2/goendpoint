package reflects

import (
	"reflect"
	"strconv"
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
	switch {
	case t.Kind() == reflect.Ptr:
		return "*" + RepresentType(t.Elem())
	case t.Kind() == reflect.Map:
		return "map[" + RepresentType(t.Key()) + "]" + RepresentType(t.Elem())
	case t.Kind() == reflect.Slice:
		return "[]" + RepresentType(t.Elem())
	case t.Kind() == reflect.Array:
		return "[" + strconv.FormatInt(int64(t.Len()), 10) + "]" + RepresentType(t.Elem())
	}

	if t.PkgPath() == "" {
		return t.String()
	}
	return t.PkgPath() + "." + t.Name()
}

func IsEmptyInterface(t reflect.Type) bool {
	return RepresentType(t) == "interface {}"
}
