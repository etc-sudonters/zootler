package vm

import (
	"context"
	"errors"
	"fmt"
	"sudonters/zootler/internal/rules/bytecode"
	"sudonters/zootler/internal/slipup"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

var ErrStackOverflow = errors.New("stackoverflow")

const maxStack = 256

func Evaluate(ctx context.Context, chunk *bytecode.Chunk) (*Execution, error) {
	vm := &Execution{
		Chunk: chunk,
		pc:    0,
		stack: stack.Make[bytecode.Value](0, 256),
	}

	err := vm.Run(ctx)
	return vm, err
}

type Execution struct {
	Chunk *bytecode.Chunk
	pc    int
	stack *stack.S[bytecode.Value]
}

func (v *Execution) Result() bytecode.Value {
	if v.stack.Len() == 0 {
		return bytecode.NullValue()
	}

	return peekStack(v.stack)
}

func peekStack(s *stack.S[bytecode.Value]) bytecode.Value {
	return (*s)[s.Len()-1]
}

func (v *Execution) Run(_ context.Context) error {
	v.pc = 0
	size := v.Chunk.Len()
	pos := -1

loop:
	for v.pc < size {
		if pos == v.pc {
			panic(fmt.Errorf(
				"'%s' did not increment program counter",
				bytecode.Bytecode(v.Chunk.Ops[pos]),
			))
		}
		pos = v.pc
		switch op := bytecode.Bytecode(v.Chunk.Ops[v.pc]); op {
		case bytecode.OP_NOP:
			v.pc++
			break
		case bytecode.OP_RETURN:
			v.pc++
			break loop
		case bytecode.OP_CONST:
			idx := v.Chunk.Ops[v.pc+1]
			v.pushStack(v.Chunk.Constants[idx])
			v.pc += 2
			break
		case bytecode.OP_EQ:
			v.pushStack(bytecode.ValueFromBool(v.popStack().Eq(v.popStack())))
			v.pc++
		case bytecode.OP_NEQ:
			v.pushStack(bytecode.ValueFromBool(!v.popStack().Eq(v.popStack())))
			v.pc++
		case bytecode.OP_LT:
			// care about order
			b := v.popStack()
			a := v.popStack()
			v.pushStack(bytecode.ValueFromBool(a.Lt(b)))
			v.pc++
		case bytecode.OP_DEBUG_STACK:
			dumpStack(v.stack)
			v.pc++
			break
		case bytecode.OP_JUMP_FALSE:
			test := v.popStack().Truthy()
			if !test {
				offset := v.Chunk.ReadU16(bytecode.PC(v.pc + 1))
				v.pc += int(offset)
				break
			}
			v.pc += 3
			break
		case bytecode.OP_JUMP:
			offset := v.Chunk.ReadU16(bytecode.PC(v.pc + 1))
			v.pc += int(offset)
			break
		case bytecode.OP_DUP:
			dup := v.popStack()
			v.pushStack(dup)
			v.pushStack(dup)
			v.pc++
			break
		case bytecode.OP_ROTATE2:
			pop1 := v.popStack()
			pop2 := v.popStack()
			v.pushStack(pop1)
			v.pushStack(pop2)
			v.pc++
			break
		case bytecode.OP_AND:
			lhs, rhs := v.popStack(), v.popStack()
			v.pushStack(bytecode.ValueFromBool(lhs.Truthy() && rhs.Truthy()))
			v.pc++
			break
		case bytecode.OP_OR:
			lhs, rhs := v.popStack(), v.popStack()
			v.pushStack(bytecode.ValueFromBool(lhs.Truthy() || rhs.Truthy()))
			v.pc++
			break
		default:
			panic(notImpl(op))
		}
	}

	return nil
}

func (v *Execution) popStack() bytecode.Value {
	val, err := v.stack.Pop()
	if err != nil {
		panic(slipup.Trace(err, "vm stack"))
	}
	return val
}

func (v *Execution) pushStack(val bytecode.Value) error {
	if v.stack.Len() >= maxStack {
		return ErrStackOverflow
	}
	v.stack.Push(val)
	return nil
}

func dumpStack(s *stack.S[bytecode.Value]) {
	fmt.Printf("[Stack: \n")
	for i := s.Len() - 1; i >= 0; i-- {
		fmt.Printf("  %#v,\n", (*s)[i])
	}
	fmt.Printf("]\n")
}

func notImpl(op bytecode.Bytecode) error {
	return fmt.Errorf("%s not implemented", op)
}
