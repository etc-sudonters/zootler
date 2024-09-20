package compiler

import (
	"sudonters/zootler/icearrow/zasm"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

func LastMileOptimizations(st *SymbolTable, intrinsics *Intrinsics) func(CompileTree) CompileTree {
	callIntrinsics := CallIntrinsics(intrinsics, st)
	//	reductions := CompressReductions()
	//	repeatedHas := CompressRepeatedHas(st)

	return func(ct CompileTree) CompileTree {
		ct = walktree(&callIntrinsics, ct)
		// ct = walktree(&reductions, ct)
		// return walktree(&repeatedHas, ct)
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
				// either we drag the nested reduction into us, or we walks its
				// fragment to compact levels within it
				if item.Op == reduction.Op {
					for _, t := range item.Targets {
						reducing.Push(t)
					}
					break
				}
				building.Targets = append(building.Targets, walktree(tw, item))
				break
			default:
				building.Targets = append(building.Targets, item)
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
	one := st.ConstOf(zasm.Pack[uint8](1))

	walker.reduce = func(tw *treewalk, reduction Reduction) CompileTree {
		var call Invocation
		var node Reduction
		node.Op = reduction.Op
		processing := stack.From(reduction.Targets)

		switch node.Op {
		case CT_REDUCE_AND:
			call.Id = hasAll.Id
			break
		case CT_REDUCE_OR:
			call.Id = hasAny.Id
			break
		default:
			panic("unreachable")
		}

		for processing.Len() > 0 {
			item, _ := processing.Pop()
			switch item := item.(type) {
			case Reduction:
				node.Targets = append(node.Targets, walktree(tw, item))
				break
			case Invocation:
				if item.Id != has.Id {
					node.Targets = append(node.Targets, item)
				}
				qty, isLoad := item.Args[1].(Load)
				if !isLoad || qty.Kind != CT_LOAD_CONST || qty.Id != one.Id {
					return node
				}
				call.Args = append(call.Args, item.Args[0].(Load))
				break
			default:
				node.Targets = append(node.Targets, item)
				break
			}
		}
		node.Targets = append(node.Targets, call)
		return node

	}
	return walker
}
