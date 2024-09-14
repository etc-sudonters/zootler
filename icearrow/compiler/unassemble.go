package compiler

import (
	"sudonters/zootler/icearrow/zasm"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

type SymbolTable struct{}
type SymbolData struct {
	Kind LoadKind
}

func (st *SymbolTable) Const(handle uint32) (SymbolData, bool) {
	var data SymbolData
	return data, false
}

func (st *SymbolTable) Symbol(handle uint32) (SymbolData, bool) {
	var data SymbolData
	return data, false
}

// Loads an assembly into a graph structure so we can plug final values such as
// settings before the final code gen
func Unassemble(asm zasm.Unit, st *SymbolTable) (CompileTree, error) {
	var graph CompileTree
	instructions := stack.From(asm.I)
	graph, instructions = unassemble(instructions, st)
	if instructions.Len() > 0 {
		return graph, slipup.Createf("failed to read all instructions into compile tree")
	}
	return graph, nil
}

func unassemble(instructions *stack.S[zasm.Instruction], st *SymbolTable) (CompileTree, *stack.S[zasm.Instruction]) {
	instruction, _ := instructions.Pop()
	op, payload := instruction.OpAndPayload()
	switch op {
	// Values
	case zasm.OP_LOAD_CONST:
		handle := zasm.DecodeU24(payload)
		constant, _ := st.Symbol(handle)
		return Load{Id: handle, Kind: constant.Kind}, instructions
	case zasm.OP_LOAD_IDENT:
		handle := zasm.DecodeU24(payload)
		sym, _ := st.Symbol(handle)
		return Load{Id: handle, Kind: sym.Kind}, instructions
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
		// Productions
	case zasm.OP_CMP_EQ:
		return production(CT_PRODUCE_EQ, instructions, st)
	case zasm.OP_CMP_NQ:
		return production(CT_PRODUCE_NQ, instructions, st)
	case zasm.OP_CMP_LT:
		return production(CT_PRODUCE_LT, instructions, st)
		// Reductions
	case zasm.OP_BOOL_AND:
		return reduction(CT_REDUCE_AND, instructions, st)
	case zasm.OP_BOOL_OR:
		return reduction(CT_REDUCE_OR, instructions, st)
	case zasm.OP_BOOL_NEGATE:
		var target CompileTree
		target, instructions = unassemble(instructions, st)
		return Inversion{target}, instructions
		// Invocations
	case zasm.OP_CALL_0:
		return Invocation{Id: zasm.DecodeU24(payload)}, instructions
	case zasm.OP_CALL_1:
		var invoke Invocation
		invoke.Id = zasm.DecodeU24(payload)
		invoke.Args = []CompileTree{nil}
		invoke.Args[0], instructions = unassemble(instructions, st)
		return invoke, instructions
	case zasm.OP_CALL_2:
		var invoke Invocation
		invoke.Id = zasm.DecodeU24(payload)
		invoke.Args = []CompileTree{nil, nil}
		// the assembly builds the stack [arg 2] [arg 1] however, we're reading
		// the instructions backwards so it looks [arg 1] [arg 2] to us
		invoke.Args[0], instructions = unassemble(instructions, st)
		invoke.Args[1], instructions = unassemble(instructions, st)
		return invoke, instructions
	default:
		panic("unknown op")
	}
}

func production(kind Producer, instructions *stack.S[zasm.Instruction], st *SymbolTable) (Production, *stack.S[zasm.Instruction]) {
	var lhs, rhs CompileTree
	lhs, instructions = unassemble(instructions, st)
	rhs, instructions = unassemble(instructions, st)
	return Production{Op: kind, Targets: []CompileTree{lhs, rhs}}, instructions
}

func reduction(kind Reducer, instructions *stack.S[zasm.Instruction], st *SymbolTable) (Reduction, *stack.S[zasm.Instruction]) {
	var lhs, rhs CompileTree
	lhs, instructions = unassemble(instructions, st)
	rhs, instructions = unassemble(instructions, st)
	return Reduction{Op: kind, Targets: []CompileTree{lhs, rhs}}, instructions
}
