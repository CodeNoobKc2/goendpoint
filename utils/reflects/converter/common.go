package converter

import (
	"encoding/json"
	"reflect"
)

var pkgPath = reflect.TypeOf(RawJSON{}).PkgPath()

// RawJSON differ value from json byte and normal byte
type RawJSON []byte

var _ Conversion = JSONToStruct{}

type JSONToStruct struct{}

func (J JSONToStruct) Accept(in, out reflect.Type) bool {
	return in.PkgPath() == pkgPath && in.Name() == "RawJSON" && out.Kind() == reflect.Struct
}

func (J JSONToStruct) ConvertValue(_ *Converter, in, out reflect.Value) (err error) {
	s := reflect.New(out.Type())
	if err := json.Unmarshal(in.Bytes(), s.Interface()); err != nil {
		return err
	}
	out.Set(s.Elem())
	return nil
}

// RawXML differ value from raw xml byte and normal byte
type RawXML []byte
