package converter

// Atomic tells what kind of conversion it can perform (typeA -> typeB) and perform them
type Atomic interface {
	// Accept return what kind of type accept
	Accept() string
	// Produce return what kind of type produce
	Produce() string
	// Convert perform conversion and raise corresponding error if occurred
	Convert(in interface{}) (out interface{}, err error)
}

var _ Atomic = atom{}

type atom struct {
	accept  string
	produce string
	convert func(in interface{}) (out interface{}, err error)
}

func (a atom) Accept() string {
	return a.accept
}

func (a atom) Produce() string {
	return a.produce
}

func (a atom) Convert(in interface{}) (out interface{}, err error) {
	return a.convert(in)
}
