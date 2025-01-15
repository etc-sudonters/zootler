package vm

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
)

func GlobalNames() []string {
	return globalNames[:]
}

var globalNames = []string{
	"Fire",
	"Forest",
	"Light",
	"Shadow",
	"Spirit",
	"Water",
	"adult",
	"age",
	"both",
	"either",
	"child",
}

var (
	ErrStackEmpty = errors.New("stack empty")
	ErrStackFull  = errors.New("stack full")
)

func newstack[T any](size int) stack[T] {
	return stack[T]{
		items: make([]T, size),
	}
}

type stack[T any] struct {
	items []T
	ptr   int
}

func (this *stack[T]) reset() {
	this.ptr = 0
}

func (this *stack[T]) slice(start, count int) []T {
	return this.items[start : start+count]
}

func (this *stack[T]) popN(n int) {
	this.ptr -= n
}

func (this *stack[T]) top() T {
	if this.ptr < 1 {
		panic(ErrStackEmpty)
	}

	return this.items[this.ptr-1]
}

func (this *stack[T]) push(item T) {
	if this.ptr == len(this.items) {
		panic(ErrStackFull)
	}

	this.items[this.ptr] = item
	this.ptr++
}

func (this *stack[T]) pop() T {
	if this.ptr == 0 {
		panic(ErrStackEmpty)
	}
	this.ptr--
	return this.items[this.ptr]
}

type VM struct {
	Objects *objects.Table
}

type execution struct {
	ip    int
	code  compiler.Bytecode
	stack stack[objects.Object]
}

func (this *execution) reset() {
	this.ip = 0
}

func (this *execution) readu8() uint8 {
	u8 := code.ReadU8(this.code.Tape[this.ip:])
	this.ip++
	return u8
}

func (this *execution) readu16() uint16 {
	u16 := code.ReadU16(this.code.Tape[this.ip:])
	this.ip += 2
	return u16
}

func (this *execution) endOfTape() int {
	return len(this.code.Tape)
}

func (this *execution) readIndex() objects.Index {
	return objects.Index(this.readu16())
}

func (this *execution) readOp() code.Op {
	return code.Op(this.readu8())
}

func (this *execution) stackargs(count int) []objects.Object {
	return this.stack.slice(this.stack.ptr-count, count)
}

func (this *VM) Execute(bytecode compiler.Bytecode) (result objects.Object, err error) {
	unit := execution{0, bytecode, newstack[objects.Object](256)}
	EOT := unit.endOfTape()

	defer func() {
		if r := recover(); r != nil {
			if thisErr, ok := r.(error); ok {
				err = thisErr
			} else if str, ok := r.(string); ok {
				err = fmt.Errorf(str)
			}

			err = fmt.Errorf("PANIC!!!! %w\n%s", err, debug.Stack())
		}
	}()

loop:
	for unit.ip < EOT {
		thisOp := unit.readOp()
		switch thisOp {
		case code.NOP:
			continue
		case code.ERR:
			err = errors.Join(errors.New("execution halted"), err)
			break loop
		case code.PUSH_T:
			unit.stack.push(objects.Boolean(true))
		case code.PUSH_F:
			unit.stack.push(objects.Boolean(false))
		case code.PUSH_CONST:
			index := unit.readIndex()
			constant := this.Objects.Constant(index)
			unit.stack.push(constant)
		case code.PUSH_TOKEN, code.PUSH_SETTING:
			ptr := unit.readIndex()
			unit.stack.push(this.Objects.Pointer(ptr))
		case code.PUSH_BUILTIN:
			index := unit.readIndex()
			builtin := this.Objects.BuiltIn(index)
			unit.stack.push(builtin)
		case code.INVERT:
			obj := unit.stack.pop()
			unit.stack.push(objects.Boolean(!this.truthy(obj)))
		case code.NEED_ALL:
			count := int(unit.readu16())
			var reduction objects.Boolean = true
			stackargs := unit.stackargs(count)
			for _, obj := range stackargs {
				if !this.truthy(obj) {
					reduction = false
					break
				}
			}
			unit.stack.popN(count)
			unit.stack.push(reduction)
		case code.NEED_ANY:
			count := int(unit.readu16())
			var reduction objects.Boolean = false
			stackargs := unit.stackargs(count)
			for _, obj := range stackargs {
				if this.truthy(obj) {
					reduction = true
					break
				}
			}
			unit.stack.popN(count)
			unit.stack.push(reduction)
		case code.CHK_QTY:
			var answer objects.Object
			ptr := unit.readIndex()
			qty := unit.readu8()
			answer, err = func([]objects.Object) (objects.Object, error) {
				return objects.Boolean(true), nil
			}([]objects.Object{
				this.Objects.Pointer(ptr), objects.Number(qty),
			})
			if err != nil {
				err = fmt.Errorf("has 0x%04x 0x%02x: %w", ptr, qty, err)
				break loop
			}
			unit.stack.push(answer)
		case code.INVOKE:
			fn := unit.stack.pop()
			var answer objects.Object
			switch fn := fn.(type) {
			case *objects.BuiltInFunction:
				count := int(unit.readu16())
				if fn.Params > -1 && count != fn.Params {
					err = fmt.Errorf("%q expects %d arguments, received %d", fn.Name, fn.Params, count)
					break loop
				}
				args := unit.stackargs(count)
				answer, err = fn.Fn(args)
				if err != nil {
					err = fmt.Errorf("%q: %w", fn.Name, err)
					break loop
				}
				unit.stack.popN(count)
			default:
				err = fmt.Errorf("cannot call %s", fn)
				break loop
			}
			if answer != nil {
				unit.stack.push(answer)
			}
		case code.CMP_EQ, code.CMP_NQ, code.CMP_LT:
			err = fmt.Errorf("runtime comparison not implemented")
			break loop
		case code.CHK_ALL, code.CHK_ANY, code.IS_CHILD, code.IS_ADULT:
			err = fmt.Errorf("fastop 0x%02x not implemented", thisOp)
			break loop
		default:
			err = fmt.Errorf("unrecognized op: 0x%02x", thisOp)
			break loop
		}
	}

	if err == nil && unit.stack.ptr > 0 {
		result = unit.stack.pop()
	}
	return
}

func (this *VM) truthy(obj objects.Object) bool {
	switch obj := obj.(type) {
	case objects.Boolean:
		return bool(obj)
	case objects.Number:
		return obj != 0
	case objects.String:
		return obj != ""
	default:
		return true
	}
}
