package objects

type Callable = func(...Object) (Object, error)

type Kind string
type Boolean bool
type String string
type Number float64
type BuiltInFunc struct {
	Name string
	Func Callable
}

const (
	_        Kind = ""
	BOOLEAN       = "BOOLEAN"
	BUILT_IN      = "BUILT_IN"
	NUMBER        = "NUMBER"
	STRING        = "STRING"
)

type Object interface {
	Kind() Kind
}

func (this String) Kind() Kind {
	return STRING
}

func (this Number) Kind() Kind {
	return NUMBER
}

func (this Boolean) Kind() Kind {
	return BOOLEAN
}

func (this BuiltInFunc) Kind() Kind {
	return BUILT_IN
}
