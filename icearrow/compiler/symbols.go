package compiler

type SymbolTable struct {
	byname map[string]int
	ident  []Symbol

	consts  []Const
	strings []String
}

type Symbol struct {
	Id    uint32
	Flags uint32
	Name  string
}

type Const struct {
	Id    uint32
	Value float64
}

type String struct {
	Offset uint16
	Len    uint8
}
