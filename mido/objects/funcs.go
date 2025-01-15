package objects

type BuiltInFunctionDef struct {
	Name   string
	Params int
}
type BuiltInFunction func(*Table, []Object) (Object, error)
