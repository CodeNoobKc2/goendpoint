package reflects

import (
	"fmt"
	"reflect"
)

func ShouldBe(v reflect.Value, kind reflect.Kind) error {
	if v.Kind() != kind {
		return fmt.Errorf("kind should be %v", kind.String())
	}
	return nil
}

func Underlying(v reflect.Value, recursive bool) reflect.Value {
	return v
}

// IsNumber return number if true
func IsNumber(t reflect.Type) bool {
	return t.Kind() == reflect.Int || t.Kind() == reflect.Int8 || t.Kind() == reflect.Int32 || t.Kind() == reflect.Int64 ||
		t.Kind() == reflect.Uint || t.Kind() == reflect.Uint8 || t.Kind() == reflect.Uint16 || t.Kind() == reflect.Uint32 || t.Kind() == reflect.Uint64 ||
		t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
}

// IsString return string
func IsString(t reflect.Type) bool {
	return t.Kind() == reflect.String
}

// IsObject return ture on interface or struct
func IsObject(t reflect.Type) bool {
	return t.Kind() == reflect.Interface || t.Kind() == reflect.Struct
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
