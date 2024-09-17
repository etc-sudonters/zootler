package compiler

type RuleCompiler struct {
	Symbols *SymbolTable
}

func (rc *RuleCompiler) Compile(tree CompileTree) tape {

	tape := new(tape)
	return *tape
}
