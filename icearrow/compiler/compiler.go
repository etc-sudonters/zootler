package compiler

type RuleCompiler struct {
	Symbols *SymbolTable
}

func (rc *RuleCompiler) Compile(fragment CompileTree) tape {
	var tw treewalk
	tape := new(tape)
	tw.immediate = func(t *treewalk, i Immediate) CompileTree {
		switch i.Kind {
		case CT_IMMED_FALSE:
			tape.writeLoadFalse()
			break
		case CT_IMMED_TRUE:
			tape.writeLoadTrue()
			break
		case CT_IMMED_U8:
			tape.writeLoadImmediateU8(i.Value.(uint8))
			break
		case CT_IMMED_U16:
			tape.writeLoadImmediateU16(i.Value.(uint16))
			break
		default:
			panic("unreachable")
		}
		return i
	}

	tw.invoke = func(t *treewalk, i Invocation) CompileTree {
		for idx := range i.Args {
			walktree(t, i.Args[idx])
		}
		tape.writeCall(uint16(i.Id), uint8(len(i.Args)))
		return i
	}

	tw.load = func(t *treewalk, l Load) CompileTree {
		switch l.Kind {
		case CT_LOAD_CONST:
			tape.writeLoadConst(uint16(l.Id))
			break
		case CT_LOAD_IDENT:
			tape.writeLoadSymbol(uint16(l.Id))
			break
		case CT_LOAD_STR:
			tape.writeLoadString(uint16(l.Id))
			break
		}
		return l
	}

	tw.reduce = func(t *treewalk, r Reduction) CompileTree {
		for i := range r.Targets {
			walktree(t, r.Targets[i])
		}
		switch r.Op {
		case CT_REDUCE_AND:
			tape.writeReduceAll(uint8(len(r.Targets)))
			break
		case CT_REDUCE_OR:
			tape.writeReduceAny(uint8(len(r.Targets)))
			break
		default:
			panic("unknown reducer")

		}
		return r
	}

	walktree(&tw, fragment)
	return *tape
}
