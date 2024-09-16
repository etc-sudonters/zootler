package compiler

type RuleCompiler struct{}

func (rc *RuleCompiler) Compile(tree CompileTree) tape {
	tape := new(tape)
	return *tape
}
