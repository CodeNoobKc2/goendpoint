package sets

type GenericHashSet map[interface{}]struct{}

func NewGenericHashSet() *GenericHashSet {
	m := make(GenericHashSet)
	return &m
}

func (s *GenericHashSet) Has(i interface{}) bool {
	_, ok := (*s)[i]
	return ok
}

func (s *GenericHashSet) Insert(i interface{}) {
	(*s)[i] = struct{}{}
}

func (s GenericHashSet) Len() int {
	return len(s)
}

func (s GenericHashSet) UnsortedList() []interface{} {
	ret := make([]interface{}, s.Len())
	idx := 0
	for item, _ := range s {
		ret[idx] = item
		idx++
	}
	return ret
}
