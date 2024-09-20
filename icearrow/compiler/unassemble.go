package compiler

import (
	"sudonters/zootler/icearrow/zasm"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

// Loads an assembly into a graph structure so we can plug final values such as
// settings before the final code gen
func Unassemble(asm *zasm.Unit, st *SymbolTable) (CompileTree, error) {
	var graph CompileTree
	instructions := stack.From(asm.I)
	graph, instructions = unassemble(instructions, st)
	if instructions.Len() > 0 {
		return graph, slipup.Createf("failed to read all instructions into compile tree")
	}
	return graph, nil
}

var loadSymKinds = map[zasm.Op]SymbolKind{
	zasm.OP_LOAD_TOK: SYM_KIND_TOKEN,
	zasm.OP_LOAD_SYM: SYM_KIND_SYMBOL,
	zasm.OP_LOAD_VAR: SYM_KIND_VAR,
	zasm.OP_LOAD_TRK: SYM_KIND_TRICK,
	zasm.OP_LOAD_SET: SYM_KIND_SETTING,
}

func unassemble(instructions *stack.S[zasm.Instruction], st *SymbolTable) (CompileTree, *stack.S[zasm.Instruction]) {
	instruction, _ := instructions.Pop()
	op, payload := instruction.OpAndPayload()
	switch op {
	// Values
	case zasm.OP_LOAD_CONST:
		handle := zasm.DecodeU24(payload)
		return Load{Id: handle, Kind: CT_LOAD_CONST}, instructions
	case zasm.OP_LOAD_TOK,
		zasm.OP_LOAD_SYM,
		zasm.OP_LOAD_VAR,
		zasm.OP_LOAD_TRK,
		zasm.OP_LOAD_SET: // end cases

		handle := zasm.DecodeU24(payload)
		symbol := st.Symbol(handle)
		kind, exists := loadSymKinds[op]
		if !exists {
			panic(slipup.Createf("unknown symbol mapping for op 0x%02x", op))
		}
		symbol.Set(kind)
		return Load{Id: handle, Kind: CT_LOAD_IDENT}, instructions
	case zasm.OP_LOAD_STR:
		handle := zasm.DecodeU24(payload)
		return Load{Id: handle, Kind: CT_LOAD_STR}, instructions
	case zasm.OP_LOAD_BOOL:
		truthy := zasm.DecodeBool(payload)
		kind := CT_IMMED_FALSE
		if truthy {
			kind = CT_IMMED_TRUE
		}
		return Immediate{Value: truthy, Kind: kind}, instructions
		// Reductions
	case zasm.OP_BOOL_AND:
		return reduction(CT_REDUCE_AND, instructions, st)
	case zasm.OP_BOOL_OR:
		return reduction(CT_REDUCE_OR, instructions, st)
		// Invocations
	case zasm.OP_CALL_0:
		handle := zasm.DecodeU24(payload)
		symbol := st.Symbol(handle)
		symbol.Set(SYM_KIND_CALLABLE)
		return Invocation{Id: handle}, instructions
	case zasm.OP_CALL_1:
		var invoke Invocation
		handle := zasm.DecodeU24(payload)
		symbol := st.Symbol(handle)
		symbol.Set(SYM_KIND_CALLABLE)
		invoke.Id = handle
		invoke.Args = []CompileTree{nil}
		invoke.Args[0], instructions = unassemble(instructions, st)
		return invoke, instructions
	case zasm.OP_CALL_2:
		var invoke Invocation
		handle := zasm.DecodeU24(payload)
		symbol := st.Symbol(handle)
		symbol.Set(SYM_KIND_CALLABLE)
		invoke.Id = handle
		invoke.Args = []CompileTree{nil, nil}
		invoke.Args[1], instructions = unassemble(instructions, st)
		invoke.Args[0], instructions = unassemble(instructions, st)
		return invoke, instructions
	default:
		panic(slipup.Createf("unknown zasm op 0x%2X", op))
	}
}

func reduction(kind Reducer, instructions *stack.S[zasm.Instruction], st *SymbolTable) (Reduction, *stack.S[zasm.Instruction]) {
	var lhs, rhs CompileTree
	lhs, instructions = unassemble(instructions, st)
	rhs, instructions = unassemble(instructions, st)
	return Reduction{Op: kind, Targets: []CompileTree{lhs, rhs}}, instructions
}
