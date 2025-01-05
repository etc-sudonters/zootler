package objects

import (
	"sudonters/zootler/magicbeanvm/nan"
)

type Kind string
type Boolean bool
type Token nan.Pointer
type String []byte
type Number float64
type BuiltIn struct {
	Func      func(...Object) (Object, error)
	NumParams int
}

const (
	_        Kind = ""
	BOOLEAN       = "BOOLEAN"
	BUILT_IN      = "BUILT_IN"
	NUMBER        = "NUMBER"
	STRING        = "STRING"
	TOKEN         = "TOKEN"
)

type Object interface {
	Kind() Kind
}

func (this *Token) Kind() Kind {
	return TOKEN
}
func (this *String) Kind() Kind {
	return STRING
}
func (this *Number) Kind() Kind {
	return NUMBER
}
func (this *Boolean) Kind() Kind {
	return BOOLEAN
}

func (this *BuiltIn) Kind() Kind {
	return BUILT_IN
}
