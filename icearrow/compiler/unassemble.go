package compiler

import (
	"sudonters/zootler/icearrow/zasm"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

// Loads an assembly into a graph structure so we can plug final values such as
// settings before the final code gen
func Unassemble(asm zasm.Unit) (CompileTree, error) {
	var graph CompileTree
	instructions := stack.From(asm.I)
	graph, instructions = unassemble(instructions)
	if instructions.Len() > 0 {
		return graph, slipup.Createf("failed to read all instructions into compile tree")
	}
	return graph, nil
}

func unassemble(instructions *stack.S[zasm.Instruction]) (CompileTree, *stack.S[zasm.Instruction]) {
	instruction, _ := instructions.Pop()
	op, payload := instruction.OpAndPayload()
	switch op {
	// Values
	case zasm.OP_LOAD_CONST:
		handle := zasm.DecodeU24(payload)
		return Load{Id: handle, Kind: CT_LOAD_CONST}, instructions
	case zasm.OP_LOAD_IDENT:
		handle := zasm.DecodeU24(payload)
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
		// Productions
	case zasm.OP_CMP_EQ:
		return production(CT_PRODUCE_EQ, instructions)
	case zasm.OP_CMP_NQ:
		return production(CT_PRODUCE_NQ, instructions)
	case zasm.OP_CMP_LT:
		return production(CT_PRODUCE_LT, instructions)
		// Reductions
	case zasm.OP_BOOL_AND:
		return reduction(CT_REDUCE_AND, instructions)
	case zasm.OP_BOOL_OR:
		return reduction(CT_REDUCE_OR, instructions)
	case zasm.OP_BOOL_NEGATE:
		var target CompileTree
		target, instructions = unassemble(instructions)
		return Inversion{target}, instructions
		// Invocations
	case zasm.OP_CALL_0:
		return Invocation{Id: zasm.DecodeU24(payload)}, instructions
	case zasm.OP_CALL_1:
		var invoke Invocation
		invoke.Id = zasm.DecodeU24(payload)
		invoke.Args = []CompileTree{nil}
		invoke.Args[0], instructions = unassemble(instructions)
		return invoke, instructions
	case zasm.OP_CALL_2:
		var invoke Invocation
		invoke.Id = zasm.DecodeU24(payload)
		invoke.Args = []CompileTree{nil, nil}
		invoke.Args[1], instructions = unassemble(instructions)
		invoke.Args[0], instructions = unassemble(instructions)
		return invoke, instructions
	default:
		panic("unknown op")
	}
}

func production(kind Producer, instructions *stack.S[zasm.Instruction]) (Production, *stack.S[zasm.Instruction]) {
	var lhs, rhs CompileTree
	lhs, instructions = unassemble(instructions)
	rhs, instructions = unassemble(instructions)
	return Production{Op: kind, Targets: []CompileTree{lhs, rhs}}, instructions
}

func reduction(kind Reducer, instructions *stack.S[zasm.Instruction]) (Reduction, *stack.S[zasm.Instruction]) {
	var lhs, rhs CompileTree
	lhs, instructions = unassemble(instructions)
	rhs, instructions = unassemble(instructions)
	return Reduction{Op: kind, Targets: []CompileTree{lhs, rhs}}, instructions
}
