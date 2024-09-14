package compiler

func CompressReductions(ct CompileTree) CompileTree {
	return ct
}

func CompressRepeatedHas(ct CompileTree) CompileTree {
	return ct
}

func InlineSeedSettings(ct CompileTree, settings any) CompileTree {
	return ct
}

func LastMileOptimizations(ct CompileTree, st *SymbolTable, settings any) CompileTree {
	ct = InlineSeedSettings(ct, settings)
	ct = CompressReductions(ct)
	ct = CompressRepeatedHas(ct)
	return ct
}

type treewalk struct {
	load      func(*treewalk, Load) CompileTree
	immediate func(*treewalk, Immediate) CompileTree
	invoke    func(*treewalk, Invocation) CompileTree
	reduce    func(*treewalk, Reduction) CompileTree
	produce   func(*treewalk, Production) CompileTree
	invert    func(*treewalk, Inversion) CompileTree
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
