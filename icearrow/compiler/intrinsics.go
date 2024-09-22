package compiler

type Intrinsic func(Invocation, *Symbol, *SymbolTable) (CompileTree, error)

func NewIntrinsics() Intrinsics {
	return Intrinsics{
		funcs: make(map[uint32]Intrinsic),
	}
}

type Intrinsics struct {
	funcs map[uint32]Intrinsic
}

func (i *Intrinsics) Add(sym *Symbol, intrin Intrinsic) {
	i.funcs[sym.Id] = intrin
}

func (i *Intrinsics) For(sym *Symbol) Intrinsic {
	return i.funcs[sym.Id]
}

func CallIntrinsics(intrinsics *Intrinsics, st *SymbolTable) treewalk {
	var walker treewalk

	walker.invoke = func(tw *treewalk, fragment Invocation) CompileTree {
		symbol := st.Symbol(fragment.Id)
		intrinsic := intrinsics.For(symbol)
		if intrinsic == nil {
			return fragment
		}

		ct, err := intrinsic(fragment, symbol, st)
		if err != nil {
			panic(err)
		}
		return ct
	}

	return walker
}
