package converter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
)

var (
	errCannotConvertFromEmptyInterface = errors.New("cannot convert from empty interface")
)

// Converter a collection of various kind of Conversion as well as some common bitSize conversion(int32 -> int64)
// some str->number(floatxx,intxx,uintxx) conversion
// on top of these Conversion, Converter would also perform complex convert between some composite type if possible
// like []typA -> []typB, map[typA]typB -> map[typC]typD
type Converter struct {
}

func (c Converter) parseNumber(vin, vout reflect.Value) error {
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
		return fmt.Errorf("out kind %v is not number", vout.Type().String())
	}
	return nil
}

func (c Converter) convertNumber(vin, vout reflect.Value) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()

		converted := vin.Convert(vout.Type())
		vout.Set(converted)
	}()
	return err
}

func (c Converter) convertListLike(vin, vout reflect.Value) error {
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
		if err := c.ConvertReflectValue(in, convertedValue); err != nil {
			return err
		}
		converted.Index(i).Set(convertedValue)
	}
	vout.Set(converted)
	return nil
}

func (c Converter) convertMap(vin, vout reflect.Value) error {
	// do nothing
	if vin.IsNil() {
		return nil
	}

	converted := reflect.MakeMap(vout.Type())
	iter := vin.MapRange()
	for iter.Next() {
		convertedKey := reflect.New(vout.Type().Key()).Elem()
		key := reflects.UnderlyingInterface(iter.Key())
		if err := c.ConvertReflectValue(key, convertedKey); err != nil {
			return err
		}

		convertedValue := reflect.New(vout.Type().Elem()).Elem()
		val := reflects.UnderlyingInterface(iter.Value())
		if err := c.ConvertReflectValue(val, convertedValue); err != nil {
			return err
		}
		converted.SetMapIndex(convertedKey, convertedValue)
	}
	vout.Set(converted)
	return nil
}

// Convert out should be addressable
func (c Converter) Convert(in, out interface{}) error {
	vin := reflect.ValueOf(in)
	vout := reflect.ValueOf(out).Elem()

	return c.ConvertReflectValue(vin, vout)
}

func (c Converter) ConvertReflectValue(vin, vout reflect.Value) error {
	if reflects.IsEmptyInterface(vin.Type()) {
		return errCannotConvertFromEmptyInterface
	}

	switch {
	case vout.Kind() == reflect.Interface:
		if vin.CanConvert(vout.Type()) {
			vout.Set(vin.Convert(vout.Type()))
			return nil
		}
		return fmt.Errorf("'%v' cannot convert to '%v'", reflects.RepresentType(vin.Type()), reflects.RepresentType(vout.Type()))
		// parse number
	case vin.Kind() == reflect.String && reflects.IsNumber(vout.Type()):
		return c.parseNumber(vin, vout)
		// convert between numbers
	case reflects.IsNumber(vin.Type()) && reflects.IsNumber(vout.Type()):
		return c.convertNumber(vin, vout)
		// convert slices
	case (vin.Kind() == reflect.Slice && vout.Kind() == reflect.Slice) ||
		(vin.Kind() == reflect.Array && vout.Kind() == reflect.Slice) ||
		(vin.Kind() == reflect.Slice && vout.Kind() == reflect.Array) ||
		(vin.Kind() == reflect.Array && vout.Kind() == reflect.Array):
		return c.convertListLike(vin, vout)
		// convert map
	case vin.Kind() == reflect.Map && vout.Kind() == reflect.Map:
		return c.convertMap(vin, vout)
	}
	return nil
}
