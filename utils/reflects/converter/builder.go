package converter

// Builder is responsible for build Converter
type Builder struct {
	Conversions    []Conversion
	discardDefault bool
}

// DiscardDefault this function will discard all default conversions
func (b *Builder) DiscardDefault() *Builder {
	b.discardDefault = true
	return b
}

// AddConversion add new Conversion
func (b *Builder) AddConversion(conversion ...Conversion) *Builder {
	b.Conversions = append(b.Conversions, conversion...)
	return b
}

func (b Builder) Build() *Converter {
	conversions := b.Conversions
	if !b.discardDefault {
		conversions = append(conversions, StrToNumber{}, NumberToNumber{}, PtrToPtr{}, MapToMap{}, ListLikeToListLike{}, JSONToStruct{}, ToInterface{})
	}

	return &Converter{conversions: conversions}
}
