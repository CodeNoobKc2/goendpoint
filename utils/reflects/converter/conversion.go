package converter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
)

type Conversion interface {
	// Accept tells what kind of conversion (typeA->typeB) does this Conversion perform
	Accept(in, out reflect.Type) bool
	// ConvertValue perform conversion from typeIn to typeOut, out must be addressable
	// some complex conversion such as list->list would require other available conversions which are stored in *Converter
	ConvertValue(converter *Converter, vin, vout reflect.Value) (err error)
}

var _ Conversion = StrToNumber{}

// StrToNumber conversion from kind str to kind number
type StrToNumber struct{}

func (s StrToNumber) Accept(tin reflect.Type, tout reflect.Type) bool {
	return tin.Kind() == reflect.String && reflects.IsNumber(tout)
}

func (s StrToNumber) ConvertValue(_ *Converter, vin, vout reflect.Value) (err error) {
	size := int(vout.Type().Size()) * 8
	switch {
	case strings.HasPrefix(vout.Kind().String(), "float"):
		parsed, err := strconv.ParseFloat(vin.String(), size)
		if err != nil {
			return err
		}
		vout.SetFloat(parsed)
	case strings.HasPrefix(vout.Kind().String(), "int"):
		parsed, err := strconv.ParseInt(vin.String(), 0, size)
		if err != nil {
			return err
		}
		vout.SetInt(parsed)
	case strings.HasPrefix(vout.Kind().String(), "uint"):
		parsed, err := strconv.ParseUint(vin.String(), 0, size)
		if err != nil {
			return err
		}
		vout.SetUint(parsed)
	case strings.HasPrefix(vout.Kind().String(), "complex"):
		parsed, err := strconv.ParseComplex(vin.String(), size)
		if err != nil {
			return err
		}
		vout.SetComplex(parsed)
	default:
		return fmt.Errorf("vout kind %v is not number", vout.Type().String())
	}
	return nil
}

var _ Conversion = NumberToNumber{}

type NumberToNumber struct{}

func (n NumberToNumber) Accept(tin reflect.Type, tout reflect.Type) bool {
	return reflects.IsNumber(tin) && reflects.IsNumber(tout)
}

func (n NumberToNumber) ConvertValue(_ *Converter, vin, vout reflect.Value) (err error) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()

		converted := vin.Convert(vout.Type())
		vout.Set(converted)
	}()
	return
}

var _ Conversion = MapToMap{}

type MapToMap struct{}

func (m MapToMap) Accept(tin, tout reflect.Type) bool {
	return tin.Kind() == reflect.Map && tout.Kind() == reflect.Map
}

func (m MapToMap) ConvertValue(converter *Converter, vin, vout reflect.Value) (err error) {
	if vin.IsNil() {
		return nil
	}

	converted := reflect.MakeMap(vout.Type())
	iter := vin.MapRange()
	for iter.Next() {
		convertedKey := reflect.New(vout.Type().Key()).Elem()
		key := reflects.UnderlyingInterface(iter.Key())
		if err := converter.ConvertValue(key, convertedKey); err != nil {
			return err
		}

		convertedValue := reflect.New(vout.Type().Elem()).Elem()
		val := reflects.UnderlyingInterface(iter.Value())
		if err := converter.ConvertValue(val, convertedValue); err != nil {
			return err
		}
		converted.SetMapIndex(convertedKey, convertedValue)
	}
	vout.Set(converted)
	return nil
}

var _ Conversion = ListLikeToListLike{}

type ListLikeToListLike struct{}

func (l ListLikeToListLike) Accept(tin, tout reflect.Type) bool {
	return (tin.Kind() == reflect.Slice && tout.Kind() == reflect.Slice) ||
		(tin.Kind() == reflect.Array && tout.Kind() == reflect.Slice) ||
		(tin.Kind() == reflect.Slice && tout.Kind() == reflect.Array) ||
		(tin.Kind() == reflect.Array && tout.Kind() == reflect.Array)
}

func (l ListLikeToListLike) ConvertValue(converter *Converter, vin, vout reflect.Value) (err error) {
	if vout.Kind() == reflect.Array && vin.Len() > vout.Type().Len() {
		return fmt.Errorf("in size '%v' exceed array length '%v'", vin.Len(), vout.Type().Len())
	}

	var converted reflect.Value
	if vout.Kind() == reflect.Array {
		converted = reflect.New(vout.Type()).Elem()
	} else {
		converted = reflect.MakeSlice(vout.Type(), vin.Len(), vin.Cap())
	}

	for i := 0; i < vin.Len(); i++ {
		convertedValue := reflect.New(vout.Type().Elem()).Elem()
		in := reflects.UnderlyingInterface(vin.Index(i))
		if err := converter.ConvertValue(in, convertedValue); err != nil {
			return err
		}
		converted.Index(i).Set(convertedValue)
	}
	vout.Set(converted)
	return nil
}

var _ Conversion = ConcreteToPtr{}

type ConcreteToPtr struct{}

func (c ConcreteToPtr) Accept(in, out reflect.Type) bool {
	return in.Kind() != reflect.Ptr && out.Kind() == reflect.Ptr
}

func (c ConcreteToPtr) ConvertValue(converter *Converter, vin, vout reflect.Value) (err error) {
	// do nothing, eg: str("") -> *int(nil)
	if vin.IsZero() {
		return nil
	}

	if vout.IsNil() {
		vout.Set(reflect.New(vout.Type().Elem()))
	}

	return converter.ConvertValue(vin, vout.Elem())
}

var _ Conversion = PtrToPtr{}

type PtrToPtr struct{}

func (p PtrToPtr) Accept(in, out reflect.Type) bool {
	return in.Kind() == reflect.Ptr && out.Kind() == reflect.Ptr
}

func (p PtrToPtr) ConvertValue(converter *Converter, vin, vout reflect.Value) (err error) {
	if vin.IsNil() {
		return nil
	}

	// initialize pointer
	if vout.IsNil() {
		vout.Set(reflect.New(vout.Type().Elem()))
	}

	return converter.ConvertValue(vin.Elem(), vout.Elem())
}

var _ Conversion = ToInterface{}

type ToInterface struct{}

func (t ToInterface) Accept(in, out reflect.Type) bool {
	return out.Kind() == reflect.Interface && in.ConvertibleTo(out)
}

func (t ToInterface) ConvertValue(_ *Converter, vin, vout reflect.Value) (err error) {
	vout.Set(vin.Convert(vout.Type()))
	return nil
}
