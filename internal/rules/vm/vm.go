package vm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sudonters/zootler/internal/rules/bytecode"
	"sudonters/zootler/internal/slipup"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

var ErrStackOverflow = errors.New("stackoverflow")
var ErrOutOfBoundsCounter = errors.New("out of bound access for program counter")
var ErrUnboundName = errors.New("unbound name")

const maxStack = 256

// Evaluate is goroutine safe as long as chunk will not result in env being modified
func Evaluate(ctx context.Context, chunk *bytecode.Chunk, env *ExecutionEnvironment) (*Execution, error) {
	vm := &Execution{
		pc:     0,
		result: bytecode.NullValue(),
		Chunk:  chunk,
		Env:    env,
		stack:  stack.Make[bytecode.Value](0, 256),
	}

	err := vm.Run(ctx)
	return vm, err
}

type Execution struct {
	result bytecode.Value
	pc     uint16
	debug  bool
	Chunk  *bytecode.Chunk
	Env    *ExecutionEnvironment
	stack  *stack.S[bytecode.Value]
}

func (v *Execution) Debug() {
	v.debug = true
}

func (v *Execution) Result() bytecode.Value {
	return v.result
}

func peekStack(s *stack.S[bytecode.Value]) bytecode.Value {
	return (*s)[s.Len()-1]
}

func (v *Execution) Run(_ context.Context) error {
	var pos uint16 = math.MaxUint16
	v.pc = 0

	defer func() {
		if r := recover(); r != nil {
			fmt.Print(operationAt(v.Chunk, pos))
			dumpStack(v.stack)
			panic(r)
		}
	}()

	size := v.Chunk.Len()
	debug := v.debug

loop:
	for {
		if int(v.pc) > size {
			return slipup.TraceMsg(ErrOutOfBoundsCounter, operationAt(v.Chunk, pos))
		}
		if pos == v.pc {
			return fmt.Errorf(
				"'%s' did not increment program counter",
				bytecode.Bytecode(v.Chunk.Ops[pos]),
			)
		}
		pos = v.pc
		code := v.Chunk.Ops[v.pc]
		op := bytecode.Bytecode(code)
		if debug {
			fmt.Println(operationAt(v.Chunk, v.pc))
		}
		switch op {
		case bytecode.OP_NOP:
			v.pc++
			break
		case bytecode.OP_RETURN:
			v.pc++
			break loop
		case bytecode.OP_SET_RETURN:
			v.result = v.popStack()
			v.pc++
			break
		case bytecode.OP_LOAD_CONST:
			idx := v.Chunk.Ops[v.pc+1]
			v.pushStack(v.Chunk.Constants[idx])
			v.pc += 2
			break
		case bytecode.OP_LOAD_IDENT:
			idx := v.Chunk.Ops[v.pc+1]
			name := v.Chunk.Names[idx]
			value, found := v.Env.Lookup(name)
			if !found {
				return slipup.TraceMsg(ErrUnboundName, "%s : %s", name, operationAt(v.Chunk, v.pc))
			}
			v.pushStack(value)
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
		case bytecode.DEBUG_STACK_OP:
			if v.debug {
				dumpStack(v.stack)
			}
			v.pc++
			break
		case bytecode.OP_JUMP_FALSE:
			test := v.popStack().Truthy()
			if !test {
				dest := v.Chunk.ReadU16(bytecode.PC(v.pc + 1))
				v.pc = dest
				break
			}
			v.pc += 3
			break
		case bytecode.OP_JUMP_TRUE:
			test := v.popStack().Truthy()
			if test {
				dest := v.Chunk.ReadU16(bytecode.PC(v.pc + 1))
				v.pc = dest
				break
			}
			v.pc += 3
			break
		case bytecode.OP_JUMP:
			dest := v.Chunk.ReadU16(bytecode.PC(v.pc + 1))
			v.pc = dest
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

func operationAt(c *bytecode.Chunk, pos uint16) string {
	op := c.Ops[int(pos)]
	return fmt.Sprintf("handling: 0x%04X 0x%02X %s", pos, op, bytecode.Bytecode(op))
}
