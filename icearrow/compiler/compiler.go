package compiler

import "github.com/etc-sudonters/substrate/slipup"

type RuleCompiler struct {
	Symbols *SymbolTable
	fastOps map[uint32]IceArrowOp
}

func (rc *RuleCompiler) init() {
	if rc.fastOps != nil {
		return
	}

	fast := map[string]IceArrowOp{
		"has":       IA_HAS_QTY,
		"has_all":   IA_HAS_ALL,
		"has_any":   IA_HAS_ANY,
		"checkage":  IA_IS_ADULT,
		"hasbottle": IA_HAS_BOTTLE,
	}

	rc.fastOps = make(map[uint32]IceArrowOp, len(fast))
	for name, op := range fast {
		sym := rc.Symbols.Named(name)
		rc.fastOps[sym.Id] = op
	}
}

func (rc *RuleCompiler) Compile(fragment CompileTree) TapeWriter {
	rc.init()
	tape := new(TapeWriter)
	tw := rc.walkerForTape(tape)
	walktree(&tw, fragment)
	return *tape
}

func (rc *RuleCompiler) walkerForTape(tape *TapeWriter) treewalk {
	var tw treewalk
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
				switch arg := i.Args[1].(type) {
				case Load:
					var qty uint8
					switch arg.Kind {
					case CT_LOAD_CONST:
						c := rc.Symbols.Const(arg.Id)
						qty = uint8(c.Value)
						break
					case CT_LOAD_IDENT:
						var exists bool
						ident := rc.Symbols.Symbol(arg.Id)
						qty, exists = map[string]uint8{
							"lacstokens":           6,
							"bigpoecount":          1,
							"bridgetokens":         100,
							"ganonbosskeytokens":   100,
							"triforcegoalperworld": 16,
						}[ident.Name]
						if !exists {
							panic(slipup.Createf("trying to resolve setting %q", ident.Name))
						}
					default:
						panic("unreachable")
					}
					tape.write(op, handle[0], handle[1], qty)
					break
				case Immediate:
					tape.write(op, handle[0], handle[1], arg.Value.(uint8))
					break
				default:
					panic("unreachable")
				}
				break
			case IA_HAS_ANY, IA_HAS_ALL:
				for idx := range i.Args {
					walktree(t, i.Args[idx])
				}
				tape.writeLoadImmediateU8(uint8(len(i.Args)))
				tape.write(op)
				break
			case IA_IS_ADULT, IA_IS_CHILD:
				ageHandle := i.Args[0].(Load)
				age := rc.Symbols.String(ageHandle.Id)
				switch age.Value {
				case "child":
					tape.write(IA_IS_CHILD)
					break
				case "adult":
					tape.write(IA_IS_ADULT)
					break
				default:
					panic(slipup.Createf("unknown age %q", age.Value))
				}
				break
			case IA_HAS_BOTTLE:
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
		}
		return l
	}

	tw.reduce = func(t *treewalk, r Reduction) CompileTree {
		if len(r.Targets) == 0 {
			tape.writeLoadTrue()
			return r
		}

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

	return tw
}
