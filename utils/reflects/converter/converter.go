package converter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
)

var (
	errCannotConvertFromEmptyInterface = errors.New("cannot convert from empty interface")
)

// Converter a collection of various kind of Conversion as well as some common bitSize conversion(int32 -> int64)
// some str->number(floatxx,intxx,uintxx) conversion
// on top of these Conversions, Converter would also perform complex convert between some composite type if possible
// like []typA -> []typB, map[typA]typB -> map[typC]typD
// notice that: no deep copy semantic is guaranteed
type Converter struct {
	conversions []Conversion
}

// Convert out should be addressable
func (c *Converter) Convert(in, out interface{}) error {
	vin := reflect.ValueOf(in)
	vout := reflect.ValueOf(out).Elem()

	return c.ConvertValue(vin, vout)
}

func (c *Converter) ConvertValue(vin, vout reflect.Value) error {
	if reflects.IsEmptyInterface(vin.Type()) {
		return errCannotConvertFromEmptyInterface
	}

	if reflects.RepresentType(vin.Type()) == reflects.RepresentType(vout.Type()) {
		vout.Set(vin)
		return nil
	}

	for _, conversion := range c.conversions {
		if conversion.Accept(vin.Type(), vout.Type()) {
			return conversion.ConvertValue(c, vin, vout)
		}
	}

	return fmt.Errorf("unable to convert between '%v' '%v'", reflects.RepresentType(vin.Type()), reflects.RepresentType(vout.Type()))
}
