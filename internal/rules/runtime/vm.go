package runtime

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sudonters/zootler/internal/slipup"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

var ErrStackOverflow = errors.New("stackoverflow")
var ErrOutOfBoundsCounter = errors.New("out of bound access for program counter")
var ErrUnboundName = errors.New("unbound name")

const maxStack = 256

func CreateVM(globals *ExecutionEnvironment, heap *VmHeap) (vm VM) {
	vm.globals = globals
	vm.heap = heap
	return
}

type VM struct {
	globals *ExecutionEnvironment
	heap    *VmHeap
	debug   bool
}

func (v *VM) Debug(d bool) {
	v.debug = d
}

func CreateExecution(chunk *Chunk, env *ExecutionEnvironment) Execution {
	return Execution{
		Chunk: chunk,
		Env:   env,
		pc:    0,
		stack: stack.Make[Value](0, 16),
	}
}

func (v *VM) Run(ctx context.Context, chunk *Chunk) (Value, error) {
	execution := CreateExecution(chunk, v.globals)
	err := v.execute(ctx, &execution)
	return execution.GetResult(), err
}

func (v *VM) RunCompiledFunc(ctx context.Context, f *CompiledFuncValue, values Values) (Value, error) {
	execution := CreateExecution(f.chunk, f.env)
	err := v.execute(ctx, &execution)
	return execution.GetResult(), err
}

func (vm *VM) execute(ctx context.Context, execution *Execution) error {
	var pos uint16 = math.MaxUint16
	execution.pc = 0

	defer func() {
		if r := recover(); r != nil {
			fmt.Print(operationAt(execution.Chunk, pos))
			dumpStack(execution.stack)
			panic(r)
		}
	}()

	size := execution.Chunk.Len()
	debug := vm.debug || execution.debug

loop:
	for {
		if int(execution.pc) > size {
			return slipup.TraceMsg(ErrOutOfBoundsCounter, operationAt(execution.Chunk, pos))
		}
		if pos == execution.pc {
			return fmt.Errorf(
				"'%s' did not increment program counter",
				Bytecode(execution.Chunk.Ops[pos]),
			)
		}
		pos = execution.pc
		code := execution.Chunk.Ops[execution.pc]
		op := Bytecode(code)
		if debug {
			fmt.Println(operationAt(execution.Chunk, execution.pc))
		}
		switch op {
		case OP_NOP:
			execution.pc++
			break
		case OP_RETURN:
			execution.pc++
			break loop
		case OP_SET_RETURN:
			execution.result = execution.popStack()
			execution.pc++
			break
		case OP_LOAD_CONST:
			idx := execution.Chunk.Ops[execution.pc+1]
			execution.pushStack(execution.Chunk.Constants[idx])
			execution.pc += 2
			break
		case OP_LOAD_IDENT:
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, found := execution.Env.Lookup(name)
			if !found {
				return slipup.TraceMsg(ErrUnboundName, "%s : %s", name, operationAt(execution.Chunk, execution.pc))
			}
			execution.pushStack(value)
			execution.pc += 2
			break
		case OP_EQ:
			execution.pushStack(ValueFromBool(execution.popStack().Eq(execution.popStack())))
			execution.pc++
		case OP_NEQ:
			execution.pushStack(ValueFromBool(!execution.popStack().Eq(execution.popStack())))
			execution.pc++
		case OP_LT:
			// care about order
			b := execution.popStack()
			a := execution.popStack()
			execution.pushStack(ValueFromBool(a.Lt(b)))
			execution.pc++
		case DEBUG_STACK_OP:
			if debug {
				dumpStack(execution.stack)
			}
			execution.pc++
			break
		case OP_JUMP_FALSE:
			test := execution.popStack().Truthy()
			if !test {
				dest := execution.Chunk.ReadU16(PC(execution.pc + 1))
				execution.pc = dest
				break
			}
			execution.pc += 3
			break
		case OP_JUMP_TRUE:
			test := execution.popStack().Truthy()
			if test {
				dest := execution.Chunk.ReadU16(PC(execution.pc + 1))
				execution.pc = dest
				break
			}
			execution.pc += 3
			break
		case OP_JUMP:
			dest := execution.Chunk.ReadU16(PC(execution.pc + 1))
			execution.pc = dest
			break
		case OP_DUP:
			dup := execution.popStack()
			execution.pushStack(dup)
			execution.pushStack(dup)
			execution.pc++
			break
		case OP_ROTATE2:
			pop1 := execution.popStack()
			pop2 := execution.popStack()
			execution.pushStack(pop1)
			execution.pushStack(pop2)
			execution.pc++
			break
		case OP_AND:
			lhs, rhs := execution.popStack(), execution.popStack()
			execution.pushStack(ValueFromBool(lhs.Truthy() && rhs.Truthy()))
			execution.pc++
			break
		case OP_OR:
			lhs, rhs := execution.popStack(), execution.popStack()
			execution.pushStack(ValueFromBool(lhs.Truthy() || rhs.Truthy()))
			execution.pc++
			break
		case OP_CALL0:
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, err := vm.callFunc(ctx, name, nil)
			if err != nil {
				return err
			}
			execution.pushStack(value)
			execution.pc++
			break
		case OP_CALL1:
			arg := execution.popStack()
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, err := vm.callFunc(ctx, name, Values{arg})
			if err != nil {
				return err
			}
			execution.pushStack(value)
			execution.pc++
			break
		case OP_CALL2:
			arg1, arg2 := execution.popStack(), execution.popStack()
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, err := vm.callFunc(ctx, name, Values{arg1, arg2})
			if err != nil {
				return err
			}
			execution.pushStack(value)
			execution.pc++
			break
		default:
			panic(notImpl(op))
		}
	}

	return nil
}

func (vm *VM) callFunc(ctx context.Context, name string, values Values) (Value, error) {
	f := vm.heap.Funcs[name]
	if f == nil {
		return NullValue(), ErrUnboundName
	}
	return f.Run(ctx, vm, values)
}

type Execution struct {
	Chunk  *Chunk
	Env    *ExecutionEnvironment
	pc     uint16
	debug  bool
	stack  *stack.S[Value]
	result Value
}

func (e *Execution) Debug() {
	e.debug = true
}

func (e *Execution) GetResult() Value {
	return e.result
}

func (e *Execution) SetResult(v Value) {
	e.result = v
}

func (e *Execution) popStack() Value {
	val, err := e.stack.Pop()
	if err != nil {
		panic(slipup.Trace(err, "vm stack"))
	}
	return val
}

func (e *Execution) pushStack(val Value) error {
	if e.stack.Len() >= maxStack {
		return ErrStackOverflow
	}
	e.stack.Push(val)
	return nil
}

func dumpStack(s *stack.S[Value]) {
	fmt.Printf("[Stack: \n")
	for i := s.Len() - 1; i >= 0; i-- {
		fmt.Printf("  %#v,\n", (*s)[i])
	}
	fmt.Printf("]\n")
}

func notImpl(op Bytecode) error {
	return fmt.Errorf("%s not implemented", op)
}

func operationAt(c *Chunk, pos uint16) string {
	op := c.Ops[int(pos)]
	return fmt.Sprintf("handling: 0x%04X 0x%02X %s", pos, op, Bytecode(op))
}
