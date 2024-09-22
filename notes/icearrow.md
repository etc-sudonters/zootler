# IceArrowVM

The first generation VM for zootler. Ice Arrow handles, with some plugin
support, the entire life cycle of a rule:

1. Lexing, Parsing
2. Analysis
3. Intermediate Storage
4. Last mile compilation
5. Execution

Steps 1 through 3 take place in the [Saburo module](../carpenters/saburo), step
4 takes place in the [Shiro module](../carpenters/shiro), a dummied
implementation of 5 can be found at
[cmd/zootler/explore.go](../cmd/zootler/explore.go).

The compilation process interacts with and alters both the world graph -- by
creating additional edges and nodes from expanding at and here rules -- and the
run time entity tables -- by creating additional tokens and locations from
expanding at and here rules, and using the tables as source of truth for what
is considered a token in some cases. Consequentially, it requires both of these
to be initialized before analysis can begin. Compilation requires a settings
instance to be available.

## Lex and Parse

[Module](../icearrow/parser)

Lexing and parsing happen as a simultaneous step. Parsing is implemented as a
top down operator precedence, aka Pratt parsing[^fn1]. There's tons of examples
of Pratt parsing across the net, in particular Robert Nystrom's Crafting
Interpreters is what I learned from -- and in fact, much of IAVM follows from
my reading of that book.

The primary output of this module is a parse tree, or concrete syntax tree,
which corresponds heavily to the source code.

## Analysis

[AST](../icearrow/ast), [Analysis](../icearrow/analysis)

Analysis happens in two steps: first the parse tree is converted into an
abstract syntax tree. This sheds much of the complexity in the syntax, favoring
converting to function calls whenever reasonable. After this, several walks
are made of the tree to handle:

1. Tag identifiers with their actual types
2. Eliminating constant comparisons and branches
3. Expanding and inlining references to anything in helpers.json
4. Extract at/here rules
5. Rewrite `==`, `!=`. `<` and `not` (unary) into function calls

Extract at and here rules are replaced with queries for artificial events[^fn2],
and the new rule is held in the analysis context for later processing. 

## Intermediate Storage

[Module](../icearrow/zasm)

After analysis, the resulting AST and context are passed to the assembler which
creates a 32bit linear IR and data storage blocks. This IR represents the rules
before having any particular settings applied to them with the intention that
we will use this IR multiple times to create different seeds.

The assembly supports very few instructions:

1. Loading a symbol, a constant or string
2. Comparing two boolean values with AND and OR
3. Calling a function with 0, 1, or 2 arguments that are found on the stack

Comparisons existing only as function calls at this level, specifically as
calls like `compare_eq_setting`, etc. Identities are encoded as 24bit handles.
The assembler outputs the assembly in reverse polish notation. Details such
as encoding of U16, U24, etc can be found within the module.

## Last mile compilation

After settings become available, the assembly can be fully compiled into
IceArrowVM ops. First the compiler unassembles a rule into a compile tree. This
tree supports the assembly's three primary operations as well as an "Immediate"
node -- either a boolean, u8 or u16 that is encoded directly into the
instruction stream. This tree then has several passes run on it to perform:

1. Calling any intrinsics provided to the compiler
2. Ensure that we've properly unaliased a handful of items
3. Compress nested `AND(AND(...`/`OR(OR(...` into single `AND(...`/`OR(...`
   instructions
4. Gather neighboring `has(TOK, 1)` into `has_all`/`has_any` as appropriate
5. Replace empty or one target `AND(...` and `OR(...` as appropriate
6. Eliminate immediate bools in `AND(...` and `OR(...` with constant evaluation


Intrinsics are primarily used to handle settings values but can be used to
replace any matching call with an appropriate compiler tree fragment (that may
then be transformed as above). For example, under settings where tracking
Time of Day access is necessary, the compiler can emit a call to the
`checktod(uint8)` function instead of issuing an immediate true with settings
that do not require this tracking. Alternatively, an intrinsic can remain
unregistered and the compiler will emit the relevant op codes for invoking it.

The compiler prefers to emit dedicated instructions even if it means encoding
more information in the instruction stream. When the compiler encounters its
representation of `has(token, number)` instead of emitting several instructions:

```
0000 | 0x12 0x14 0x00   ; LOAD_SYMBOL token
0003 | 0x15 0x03        ; LOAD_IMMED_U8  3
0005 | 0x12 0x01 0x00   ; LOAD_SYMBOL has
0008 | 0x43             ; CALL_2
```

Which is 9 bytes, 4 instructions, two loads, and a call frame[^fn3], the
compiler will instead emit:

```
0000 | 0x69 0x14 0x00 0x03 ; IA_HAS_QTY (token lo bits) (token hi bits) (qty)
```

This results in complexity when decoding the instruction stream since we need
to know which instructions take arguments on the stack and which have them
encoded in the stream. I've contemplated using 32bit op codes for the VM as
well but I've stubbornly stuck to 8bit codes[^fn4].

The compiler spits out a structure called a "Tape" that has the operations
stored in a slice. This structure will be expanded to hold additional
information in the future hence not just handing out a slice of codes.

Shiro stores these tapes along with the symbol table to make them available
to the VM[^fn5].

## Execution

[Example](../cmd/zootler/explore.go)

These compiled rules can be combined with a world graph, a VM and a VM state
implementation to explore the world graph. The example uses a VMState that
always responds with true when asked a question which allows every node of the
graph to be reached on the first exploration and ergo every edge rule is
executed by the VM.

The VM is designed that a single instance can be instantiated and then reused
for all explorations. It currently exists as a zero-width type but will need to
support function calls for inquiries that will not have dedicated instructions,
e.g. `has_all_shuffled_notes_for(token)` or `has_hearts(float64)`[^fn6]:

```go
type VM struct {
    builtIns VMBuiltIns
}

type VMBuiltIns interface {
    Call(VMState, *compiler.Symbol, nan.PackedValue...) (bool, error)
}

func(vm *VM) Execute(*compiler.Tape, VMState, *compiler.SymbolTable) Execution {
    op := // ...
    switch op {
        // ...
        case CALL_0:
            sym := stk.Pop()
            result, err := vm.builtins.Call(sym)
            stk.Push(nan.PackBool(result))
            break
        case CALL_1:
            sym := stk.Pop()
            arg := stk.Pop()
            result, err := vm.builtins.Call(sym, arg)
            stk.Push(nan.PackBool(result))
            break
        case CALL_2:
            sym := stk.Pop()
            arg2 := stk.Pop()
            arg1 := stk.Pop()
            result, err := vm.builtins.Call(sym, arg1, arg2)
            stk.Push(nan.PackBool(result))
            break
        // ...
    }
}
```

---

[ln1]: https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/

[^fn1]: Nystrom covers [Pratt parsing in a write up][ln1]. This post uses Java,
    Compiling Interpreters offers a C implementation, and for a Go
    implementation other than the one used here, Writing Interpreters In Go by
    Thorsten Ball also uses this strategy. It's worth comparing at least two
    implementations if you haven't encountered this technique before.

[^fn2]: The short of these is that they're used to replace expensive searches,
    instead the run time creates this token at some location with an arbitrary
    rule associated with it. Then the run time can inquire if we've managed to
    grab that specific token rather than engaging in a graph walk.

[^fn3]: Which ice arrow doesn't even support at time of writing. You don't have
    to worry about slow function calls if you don't have any. :shrug:

[^fn4]: being able to say "bytecode" and have it be accurate is kind of
    appealing to my inner dork

[^fn5]: For now, I actually really don't like handing out the symbol table like
    this but its convenient

[^fn6]: Which will probably cause `IA_HAS_BOTTLE` to be demoted to a `CALL_0`,
    maybe the two age verification instructions as well.
