package compiler

import "github.com/etc-sudonters/substrate/skelly/stack"

type CompilerFuncs interface {
	DungeonHasShortcuts(string) bool
	IsTrialSkipped(string) bool
	LoadSetting(string) bool
	LoadSetting2(string, string) bool
	AtTimeOfDay(string) CompileTree
	CheckStartCond(string) bool
	CompareToSetting(string, CompileTree) bool
}

func CompressReductions(ct CompileTree) CompileTree {
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

	return walktree(&walker, ct)
}

func CompressRepeatedHas(ct CompileTree, st *SymbolTable) CompileTree {
	var walker treewalk
	hasAllSymbol, hasAllExists := Symbol{}, false //st.SymbolFor("has_all")
	hasAnySymbol, hasAnyExists := Symbol{}, false //st.SymbolFor("has_any")
	hasSymbol, hasExists := Symbol{}, true        //st.SymbolFor("has")

	if !hasExists {
		panic("has symbol is not bound")
	}

	if !hasAllExists || !hasAnyExists {
		return ct
	}

	walker.reduce = func(tw *treewalk, reduction Reduction) CompileTree {
		var call Invocation
		var node Reduction
		node.Op = reduction.Op
		processing := stack.From(reduction.Targets)

		switch node.Op {
		case CT_REDUCE_AND:
			call.Id = hasAllSymbol.Id
			break
		case CT_REDUCE_OR:
			call.Id = hasAnySymbol.Id
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
				if item.Id != hasSymbol.Id {
					node.Targets = append(node.Targets, item)
				}
				call.Vargs = append(call.Vargs, item.Args[0].(Load))
				break
			default:
				node.Targets = append(node.Targets, item)
				break
			}
		}

		node.Targets = append(node.Targets, call)
		return node

	}
	return walktree(&walker, ct)
}

func InlineSeedSettings(ct CompileTree, st *SymbolTable) CompileTree {
	return ct
}

func LastMileOptimizations(ct CompileTree, st *SymbolTable) CompileTree {
	ct = InlineSeedSettings(ct, st)
	ct = CompressReductions(ct)
	ct = CompressRepeatedHas(ct, st)
	return ct
}

type treewalk struct {
	immediate func(*treewalk, Immediate) CompileTree
	invert    func(*treewalk, Inversion) CompileTree
	invoke    func(*treewalk, Invocation) CompileTree
	load      func(*treewalk, Load) CompileTree
	produce   func(*treewalk, Production) CompileTree
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
	case Production:
		if tw.produce != nil {
			return tw.produce(tw, ct)
		}
		var p Production
		p.Op = ct.Op
		p.Targets = make([]CompileTree, len(ct.Targets))
		for i, trg := range ct.Targets {
			p.Targets[i] = walktree(tw, trg)
		}
		return p
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
	case Inversion:
		if tw.invert != nil {
			return tw.invert(tw, ct)
		}
		var i Inversion
		i.Target = walktree(tw, ct.Target)
		return i
	default:
		panic("unreachable")
	}
}
