package compiler

import (
	"runtime"
	"slices"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

func LastMileOptimizations(st *SymbolTable, intrinsics *Intrinsics) func(CompileTree) CompileTree {
	callIntrinsics := CallIntrinsics(intrinsics, st)
	reductions := CompressReductions()
	repeatedHas := CompressRepeatedHas(st)

	runtime.KeepAlive(reductions)
	runtime.KeepAlive(repeatedHas)
	return func(ct CompileTree) CompileTree {
		ct = walktree(&callIntrinsics, ct)
		ct = walktree(&reductions, ct)
		ct = walktree(&repeatedHas, ct)
		ct = walktree(&reductions, ct)
		ct = walktree(&repeatedHas, ct)
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
				} else if sym.Id == has.Id && trgt.Args[1].(Load).Kind == CT_LOAD_CONST {
					item := trgt.Args[0].(Load)
					qty := st.Const(trgt.Args[1].(Load).Id)
					if qty.Value == 1 {
						haser[sym.Id] = append(collected, item)
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
