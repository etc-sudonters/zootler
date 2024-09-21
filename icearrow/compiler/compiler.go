package compiler

type RuleCompiler struct {
	Symbols *SymbolTable
	fastOps map[uint32]IceArrowOp
}

func (rc *RuleCompiler) init() {
	if rc.fastOps != nil {
		return
	}

	fast := map[string]IceArrowOp{
		"has":     IA_HAS_QTY,
		"has_all": IA_HAS_ALL,
		"has_any": IA_HAS_ANY,
	}

	rc.fastOps = make(map[uint32]IceArrowOp, len(fast))
	for name, op := range fast {
		sym := rc.Symbols.Named(name)
		rc.fastOps[sym.Id] = op
	}
}

func (rc *RuleCompiler) Compile(fragment CompileTree) Tape {
	var tw treewalk
	tape := new(Tape)
	rc.init()

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
		sym := rc.Symbols.Symbol(i.Id)
		if op, exists := rc.fastOps[sym.Id]; exists {
			switch op {
			case IA_HAS_QTY:
				handle := encodeU16(uint16(i.Args[0].(Load).Id))

				qLoad := i.Args[1].(Load)
				if qLoad.Kind != CT_LOAD_CONST {
					panic("TODO: resolve this setting value")
				}

				qty := rc.Symbols.Const(qLoad.Id)

				tape.write(op, handle[0], handle[1], uint8(qty.Value))
				break
			case IA_HAS_ANY, IA_HAS_ALL:
				for idx := range i.Args {
					walktree(t, i.Args[idx])
				}
				tape.writeLoadImmediateU8(uint8(len(i.Args)))
				tape.write(op)
				break
			}
			return i
		}

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
