package runtime

import (
	"context"
	"errors"
	"fmt"
	"math"
	"github.com/etc-sudonters/substrate/slipup"

	"github.com/etc-sudonters/substrate/skelly/stack"
)

var ErrStackOverflow = errors.New("stackoverflow")
var ErrOutOfBoundsCounter = errors.New("out of bound access for program counter")
var ErrUnboundName = errors.New("unbound name")

const maxStack = 256

func CreateVM(globals *ExecutionEnvironment, mem *FuncNamespace) (vm VM) {
	vm.globals = globals
	vm.mem = mem
	return
}

type VM struct {
	globals *ExecutionEnvironment
	mem     *FuncNamespace
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

// TODO instead of recursing (indirectly), should look at managing via call stack
func (v *VM) RunCompiledFunc(ctx context.Context, f *CompiledFunc, values Values) (Value, error) {
	// write first N name+values into environment
	scope := f.env.ChildScope()
	for i := range values {
		scope.Set(f.chunk.Names[i], values[i])
	}

	execution := CreateExecution(f.chunk, scope)
	execution.debug = v.debug || execution.debug
	err := v.execute(ctx, &execution)
	return execution.GetResult(), err
}

func (vm *VM) execute(ctx context.Context, execution *Execution) error {
	if execution.Chunk.Len() == 0 {
		return slipup.Createf("empty program passed")
	}

	var pos uint16 = math.MaxUint16
	execution.pc = 0

	defer func() {
		if r := recover(); r != nil {
			fmt.Print(operationAt(execution.Chunk, pos))
			dumpStack(execution.stack)
			panic(r)
		}
	}()

	debug := vm.debug || execution.debug

loop:
	for {
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
			execution.loadConst()
			execution.pc += 2
			break
		case OP_LOAD_IDENT:
			loadErr := execution.loadName()
			if loadErr != nil {
				return loadErr
			}
			execution.pc += 2
			break
		case OP_EQ:
			execution.pushStack(ValueFromBool(execution.popStack().Eq(execution.popStack())))
			execution.pc++
		case OP_NEQ:
			execution.pushStack(ValueFromBool(!execution.popStack().Eq(execution.popStack())))
			execution.pc++
		case OP_LT:
			b := execution.popStack()
			a := execution.popStack()
			val := ValueFromBool(a.Lt(b))
			execution.pushStack(val)
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
				return slipup.Describef(err, operationAt(execution.Chunk, pos))
			}
			execution.pushStack(value)
			execution.pc += 2
			break
		case OP_CALL1:
			arg := execution.popStack()
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, err := vm.callFunc(ctx, name, Values{arg})
			if err != nil {
				return slipup.Describef(err, operationAt(execution.Chunk, pos))
			}
			execution.pushStack(value)
			execution.pc += 2
			break
		case OP_CALL2:
			arg2, arg1 := execution.popStack(), execution.popStack()
			idx := execution.Chunk.Ops[execution.pc+1]
			name := execution.Chunk.Names[idx]
			value, err := vm.callFunc(ctx, name, Values{arg1, arg2})
			if err != nil {
				return slipup.Describef(err, operationAt(execution.Chunk, pos))
			}
			execution.pushStack(value)
			execution.pc += 2
			break
		default:
			panic(notImpl(op))
		}
	}

	return nil
}

func (vm *VM) callFunc(ctx context.Context, name string, values Values) (Value, error) {
	f := vm.mem.funcs[name]
	if f == nil {
		return NullValue(), slipup.Describef(ErrUnboundName, "function '%s'", name)
	}
	value, err := f.Run(ctx, vm, values)
	if err != nil {
		err = slipup.Describef(err, "function call '%s'", name)
	}
	return value, err
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

func (e *Execution) loadConst() {
	e.pushStack(e.Chunk.GetConstAt(PC(e.pc + 1)))
}

func (e *Execution) loadName() error {
	idx := e.Chunk.Ops[e.pc+1]
	name := e.Chunk.Names[idx]
	value, found := e.Env.Lookup(name)
	if !found {
		return slipup.Describef(ErrUnboundName, "%s : %s", name, operationAt(e.Chunk, e.pc))
	}
	e.pushStack(value)
	return nil
}

func (e *Execution) popStack() Value {
	val, err := e.stack.Pop()
	if err != nil {
		panic(slipup.Describe(err, "vm stack"))
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

func (e *Execution) currentOpDisplay() string {
	return operationAt(e.Chunk, e.pc)
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
