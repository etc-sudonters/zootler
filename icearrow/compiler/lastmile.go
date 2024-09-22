package compiler

import (
	"github.com/etc-sudonters/substrate/skelly/stack"
	"slices"
)

func LastMileOptimizations(st *SymbolTable, intrinsics *Intrinsics) func(CompileTree) CompileTree {
	unalias := unalias(st)
	callIntrinsics := CallIntrinsics(intrinsics, st)
	reductions := CompressReductions()
	repeatedHas := CompressRepeatedHas(st)
	unwrap := unwrap(st)

	return func(ct CompileTree) CompileTree {
		ct = walktree(&callIntrinsics, ct)
		ct = walktree(&unalias, ct)
		ct = walktree(&reductions, ct)
		ct = walktree(&repeatedHas, ct)
		ct = walktree(&reductions, ct)
		ct = walktree(&repeatedHas, ct)
		ct = walktree(&unwrap, ct)
		return ct
	}
}

func CompressReductions() treewalk {
	var walker treewalk
	walker.reduce = func(tw *treewalk, reduction Reduction) CompileTree {
		reducing := stack.From(reduction.Targets)
		building := Reduction{
			Op: reduction.Op,
		}

		for reducing.Len() > 0 {
			item, _ := reducing.Pop()
			switch item := item.(type) {
			case Reduction:
				if item.Op == reduction.Op {
					for _, t := range item.Targets {
						reducing.Push(t)
					}
				} else {
					building.Targets = append(building.Targets, walktree(tw, item))
				}
				continue
			case Immediate:
				if item.Kind != CT_IMMED_FALSE && item.Kind != CT_IMMED_TRUE {
					building.Targets = append(building.Targets, item)
					continue
				}

				if item.Kind == CT_IMMED_FALSE {
					if building.Op == CT_REDUCE_AND {
						return item
					}

					continue
				}

				if building.Op == CT_REDUCE_OR {
					return item
				}
				continue
			default:
				building.Targets = append(building.Targets, item)
				continue
			}
		}
		return building
	}

	return walker
}

func CompressRepeatedHas(st *SymbolTable) treewalk {
	var walker treewalk
	has := st.Named("has")
	hasAll := st.Named("has_all")
	hasAny := st.Named("has_any")

	if has == nil || hasAll == nil || hasAny == nil {
		panic("something isn't registered")
	}

	walker.reduce = func(tw *treewalk, visiting Reduction) CompileTree {
		var reducer []CompileTree
		haser := map[uint32][]CompileTree{
			has.Id:    nil,
			hasAll.Id: nil,
			hasAny.Id: nil,
		}

		for _, trgt := range visiting.Targets {
			switch trgt := trgt.(type) {
			case Invocation:
				sym := st.Symbol(trgt.Id)
				collected, exists := haser[sym.Id]
				if !exists {
					reducer = append(reducer, trgt)
					continue
				}

				if sym.Id == hasAny.Id || sym.Id == hasAll.Id {
					haser[sym.Id] = append(collected, trgt.Args...)
					continue
				} else if sym.Id == has.Id {
					item := trgt.Args[0].(Load)
					var qty uint8
					switch arg := trgt.Args[1].(type) {
					case Load:
						if arg.Kind == CT_LOAD_CONST {
							val := st.Const(arg.Id)
							qty = uint8(val.Value)
						}
						break
					case Immediate:
						qty = arg.Value.(uint8)
						break
					default:
						panic("unreachable")
					}

					if qty == 1 {
						haser[has.Id] = append(haser[has.Id], item)
						continue
					}
				}
				reducer = append(reducer, trgt)
				break
			default:
				reducer = append(reducer, walktree(tw, trgt))
				break
			}
		}

		hasAll := Invocation{
			Id:   hasAll.Id,
			Args: haser[hasAll.Id],
		}

		hasAny := Invocation{
			Id:   hasAny.Id,
			Args: haser[hasAny.Id],
		}

		if visiting.Op == CT_REDUCE_AND {
			hasAll.Args = slices.Concat(hasAll.Args, haser[has.Id])
		} else {
			hasAny.Args = slices.Concat(hasAny.Args, haser[has.Id])
		}

		if len(hasAny.Args) > 0 {
			reducer = append(reducer, hasAny)
		}

		if len(hasAll.Args) > 0 {
			reducer = append(reducer, hasAll)
		}

		if len(reducer) == 1 {
			return reducer[0]
		}

		return Reduction{
			Op:      visiting.Op,
			Targets: reducer,
		}
	}
	return walker
}

func unalias(st *SymbolTable) treewalk {
	var walker treewalk

	has := st.Named("has")

	silvers := st.Named("silvergauntlets")
	golden := st.Named("goldengauntlets")
	longshot := st.Named("longshot")

	strSym := st.Named("progressivestrengthupgrade")
	hookSym := st.Named("hookshot")

	two := Immediate{Value: uint8(2), Kind: CT_IMMED_U8}
	three := Immediate{Value: uint8(3), Kind: CT_IMMED_U8}

	str := Load{Id: strSym.Id + 1, Kind: CT_LOAD_IDENT}
	hook := Load{Id: hookSym.Id + 1, Kind: CT_LOAD_IDENT}

	walker.invoke = func(_ *treewalk, invoke Invocation) CompileTree {
		sym := st.Symbol(invoke.Id)
		if sym.Id == has.Id {
			arg := st.Symbol(invoke.Args[0].(Load).Id)
			if arg.Id == silvers.Id {
				return Invocation{
					Id:   has.Id + 1,
					Args: []CompileTree{str, two},
				}
			} else if arg.Id == golden.Id {
				return Invocation{
					Id:   has.Id + 1,
					Args: []CompileTree{str, three},
				}
			} else if arg.Id == longshot.Id {
				return Invocation{
					Id:   has.Id + 1,
					Args: []CompileTree{hook, two},
				}
			}

		}

		return invoke
	}

	return walker
}

func unwrap(st *SymbolTable) treewalk {
	var walker treewalk
	has := st.Named("has")
	hasAll := st.Named("has_all")
	hasAny := st.Named("has_any")

	walker.invoke = func(t *treewalk, invoke Invocation) CompileTree {
		sym := st.Symbol(invoke.Id)
		if (sym.Id == hasAll.Id || sym.Id == hasAny.Id) && len(invoke.Args) == 1 {
			return Invocation{
				Id:   has.Id + 1,
				Args: []CompileTree{invoke.Args[0], Immediate{Value: uint8(1), Kind: CT_IMMED_U8}},
			}
		}

		return invoke
	}
	return walker
}
