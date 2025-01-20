package objects

type BuiltInFunctionDef struct {
	Name   string
	Params int
}
type BuiltInFunction func(*Table, []Object) (Object, error)
type BuiltInFunctions []BuiltInFunction

func (this BuiltInFunctions) Call(tbl *Table, callee Object, args []Object) (Object, error) {
	_, index := UnpackPtr32(callee)
	return this[index](tbl, args)
}
