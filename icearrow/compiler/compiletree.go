package compiler

type CompileTree interface {
	CTType() CTType
}

type Load struct {
	Id   uint32
	Kind LoadKind
}

type LoadKind uint8

type Immediate struct {
	Value any
	Kind  ImmediateKind
}

type ImmediateKind uint8

type Invocation struct {
	Id   uint32
	Args []CompileTree
	// compiles to LOAD ... ; LOAD IMMEDIATE len(vargs); LOAD IDENT; CALL_V;
	Vargs []CompileTree
}

// Reduces several boolean results to one boolean result via &&, ||
type Reduction struct {
	Op      Reducer
	Targets []CompileTree
}

type Reducer uint8
type CTType uint8

const (
	_ CTType = iota
	CT_TYPE_CONST
	CT_TYPE_IMMED
	CT_TYPE_SYMBOL
	CT_TYPE_INVOKE
	CT_TYPE_REDUCE
	CT_TYPE_PRDUCE

	_ LoadKind = iota
	CT_LOAD_CONST
	CT_LOAD_IDENT
	CT_LOAD_STR

	_ Reducer = iota
	CT_REDUCE_AND
	CT_REDUCE_OR

	_ ImmediateKind = iota
	CT_IMMED_TRUE
	CT_IMMED_FALSE
	CT_IMMED_U8
	CT_IMMED_U16
)

func (node Load) CTType() CTType       { return CT_TYPE_SYMBOL }
func (node Immediate) CTType() CTType  { return CT_TYPE_IMMED }
func (node Invocation) CTType() CTType { return CT_TYPE_INVOKE }
func (node Reduction) CTType() CTType  { return CT_TYPE_REDUCE }
