package converter

import "reflect"

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
func RepresentType(v reflect.Value) string {
	t := v.Type()
	return t.String()
}
