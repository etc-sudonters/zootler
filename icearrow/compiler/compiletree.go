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
}

// Reduces several boolean results to one boolean result via &&, ||
type Reduction struct {
	Op      Reducer
	Targets []CompileTree
}

type treewalk struct {
	immediate func(*treewalk, Immediate) CompileTree
	invoke    func(*treewalk, Invocation) CompileTree
	load      func(*treewalk, Load) CompileTree
	reduce    func(*treewalk, Reduction) CompileTree
}

func walktree(tw *treewalk, ct CompileTree) CompileTree {
	switch ct := ct.(type) {
	case Load:
		if tw.load != nil {
			return tw.load(tw, ct)
		}
		return ct
	case Immediate:
		if tw.immediate != nil {
			return tw.immediate(tw, ct)
		}
		return ct
	case Invocation:
		if tw.invoke != nil {
			return tw.invoke(tw, ct)
		}
		return ct
	case Reduction:
		if tw.reduce != nil {
			return tw.reduce(tw, ct)
		}
		var r Reduction
		r.Op = ct.Op
		r.Targets = make([]CompileTree, len(ct.Targets))
		for i, trg := range ct.Targets {
			r.Targets[i] = walktree(tw, trg)
		}
		return r
	default:
		panic("unreachable")
	}
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
