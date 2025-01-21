# Mido and MagicBean

A less confused compiler and VM. Mido handles compiling rules and arranging a
runtime for the rules, MagicBean creates a graph structure to navigate.

## Mido

Arranges [compiler settings and plugins](../mido/orchestrate.go).

- Scripted functions are inlined at their call site 
- Compiler functions are plugin driven and operate on AST
- Built in functions are plugin driven and operate on VM state
- Byte code generation also supports plugins



### Parsing

Lexing and parsing are handled by [internal/ruleparser](../internal/ruleparser)
and the output is lowered into [mido/ast](../mido/ast). 

- rudimentary type checking via [mido/symbols.Table](../mido/symbols/table.go)
- provides tools for hashing, rewriting, visiting, and rendering to s-expressions
- [optimization primarily targets AST](mido/optimizer)

### Compilation

AST is [compiled](../mido/compiler) into bytecode and creates an
[mido/objects.Table](../mido/objects/table.go) from a combination of the symbol
table and numeric and string constants in the AST.

- [mido/code](../mido/code) defines opcodes
- Requires TOKEN, REGION and TRANSIT symbols to be mapped into pointers externally
- Requires BUILT_IN_FUNCTION symbols to be mapped into pointers externally
- Stores pointers and constants in [mido/objects.Object](../mido/objects/nan.go)

### Execution

[mido/vm.Vm](../mido/vm/vm.go) is a stack machine that is provided an object
table and the runtime function table and executes bytecode. 

- what's to say, just a big switch statement
- each execution creates its own stack 
- only supports calls to builtin functions


