package reflects

import (
	"fmt"
	"reflect"
)

func ShouldBeKind(t reflect.Type, kind reflect.Kind) error {
	if t.Kind() != kind {
		return fmt.Errorf("kind should be %v", kind.String())
	}
	return nil
}

func Underlying(v reflect.Value, recursive bool) reflect.Value {
	return v
}

func IsList(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

// IsNumber return number if true
func IsNumber(t reflect.Type) bool {
	return t.Kind() == reflect.Int || t.Kind() == reflect.Int8 || t.Kind() == reflect.Int32 || t.Kind() == reflect.Int64 ||
		t.Kind() == reflect.Uint || t.Kind() == reflect.Uint8 || t.Kind() == reflect.Uint16 || t.Kind() == reflect.Uint32 || t.Kind() == reflect.Uint64 ||
		t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64 || t.Kind() == reflect.Complex64 || t.Kind() == reflect.Complex128
}

// IsString return string
func IsString(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

// IsObject return ture on interface or struct
func IsObject(t reflect.Type) bool {
	return t.Kind() == reflect.Interface || t.Kind() == reflect.Struct
}

func GetStructMember(v reflect.Value, name string, visitEmbedded bool) (*reflect.Value, error) {
	var embedded []reflect.Value
	for i := 0; i < v.NumField(); i++ {
		tfield := v.Type().Field(i)
		vfield := v.Field(i)

		if tfield.Name == name {
			return &vfield, nil
		}

		if tfield.Anonymous && tfield.Type.Kind() == reflect.Struct && visitEmbedded {
			embedded = append(embedded, vfield)
		}
	}

	for _, value := range embedded {
		vfield, _ := GetStructMember(value, name, true)
		if vfield != nil {
			return vfield, nil
		}
	}

	return nil, fmt.Errorf("field '%v' not found", name)
}

// VisitStruct visit struct field may be recursively
func VisitStruct(v reflect.Value, visitEmbedded bool, walk func(v reflect.Value, t reflect.StructField)) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		member := v.Field(i)
		walk(member, t.Field(i))

		if visitEmbedded {
			VisitStruct(member, true, walk)
		}
	}
}

// UnderlyingInterface recursively return the underlying value if current value's Kind() == reflect.Interface
func UnderlyingInterface(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return UnderlyingInterface(v.Elem())
	}
	return v
}
