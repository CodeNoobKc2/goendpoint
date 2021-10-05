package converter

// Converter a collection of various kind of atomic
// on top of these atomics, Converter would also perform complex convert between some composite type if possible
// eg: []typA -> []typB, map[typA]typB -> map[typC]typD
type Converter struct {
	// converters key would be typeARepr:typeBRepr
	converters map[string]Atomic
}

// Convert out should be addressable
func (c Converter) Convert(in, out interface{}) error {
	return nil
}
