package vm

import (
	"errors"
	"fmt"
	"sudonters/zootler/midologic/code"
	"sudonters/zootler/midologic/compiler"
	"sudonters/zootler/midologic/objects"
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
	objects *objects.Table
}

func New(objs *objects.Table) VM {
	var r VM
	r.objects = objs
	return r
}

type ExecutionUnit struct {
	ByteCode compiler.ByteCode
	ip       int
}

func (this *ExecutionUnit) reset() {
	this.ip = 0
}

func (this *ExecutionUnit) readu8() uint8 {
	u8 := code.ReadU8(this.ByteCode.Tape[this.ip:])
	this.ip++
	return u8
}

func (this *ExecutionUnit) readu16() uint16 {
	u16 := code.ReadU16(this.ByteCode.Tape[this.ip:])
	this.ip += 2
	return u16
}

func (this *ExecutionUnit) endOfTape() int {
	return len(this.ByteCode.Tape)
}

func (this *ExecutionUnit) readIndex() objects.Index {
	return objects.Index(this.readu16())
}

func (this *ExecutionUnit) readOp() code.Op {
	return code.Op(this.readu8())
}

func (this *VM) Execute(executing ExecutionUnit) (result objects.Object, err error) {
	stack := newstack[objects.Object](256)
	executing.reset()
	EOT := executing.endOfTape()

	defer func() {
		if r := recover(); r != nil {
			if thisErr, ok := r.(error); ok {
				err = thisErr
			} else if str, ok := r.(string); ok {
				err = fmt.Errorf(str)
			}

			err = fmt.Errorf("PANIC!!!! %w", err)
		}
	}()

loop:
	for executing.ip < EOT {
		thisOp := executing.readOp()
		switch thisOp {
		case code.NOP:
			continue
		case code.ERR:
			err = errors.Join(errors.New("execution halted"), err)
			break loop
		case code.PUSH_T:
			stack.push(objects.Boolean(true))
		case code.PUSH_F:
			stack.push(objects.Boolean(false))
		case code.PUSH_CONST:
			index := executing.readIndex()
			constant := this.objects.Constant(index)
			stack.push(constant)
		case code.PUSH_TOKEN, code.PUSH_SETTING:
			ptr := executing.readIndex()
			stack.push(this.objects.Pointer(ptr))
		case code.PUSH_BUILTIN:
			index := executing.readIndex()
			builtin := this.objects.BuiltIn(index)
			stack.push(builtin)
		case code.INVERT:
			obj := stack.pop()
			stack.push(objects.Boolean(!this.truthy(obj)))
		case code.NEED_ALL:
			count := int(executing.readu16())
			var reduction objects.Boolean = true
			stackargs := stack.slice(stack.ptr-count, count)
			for _, obj := range stackargs {
				if !this.truthy(obj) {
					reduction = false
					break
				}
			}
			stack.popN(count)
			stack.push(reduction)
		case code.NEED_ANY:
			count := int(executing.readu16())
			var reduction objects.Boolean = false
			stackargs := stack.slice(stack.ptr-count, count)
			for _, obj := range stackargs {
				if this.truthy(obj) {
					reduction = true
					break
				}
			}
			stack.ptr -= count
			stack.push(reduction)
		case code.CHK_QTY:
			var answer objects.Object
			ptr := executing.readIndex()
			qty := executing.readu8()
			answer, err = func([]objects.Object) (objects.Object, error) {
				return objects.Boolean(true), nil
			}([]objects.Object{
				this.objects.Pointer(ptr), objects.Number(qty),
			})
			if err != nil {
				err = fmt.Errorf("has 0x%04x 0x%02x: %w", ptr, qty, err)
				break loop
			}
			stack.push(answer)
		case code.INVOKE:
			fn := stack.pop()
			var answer objects.Object
			switch fn := fn.(type) {
			case *objects.BuiltInFunction:
				count := int(executing.readu16())
				if fn.Params > -1 && count != fn.Params {
					err = fmt.Errorf("%q expects %d arguments, received %d", fn.Name, fn.Params, count)
					break loop
				}
				args := stack.slice(stack.ptr-count, count)
				answer, err = fn.Fn(args)
				if err != nil {
					err = fmt.Errorf("%q: %w", fn.Name, err)
					break loop
				}
				stack.popN(count)
			default:
				err = fmt.Errorf("cannot call %s", fn)
				break loop
			}
			if answer != nil {
				stack.push(answer)
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

	if err == nil && stack.ptr > 0 {
		result = stack.pop()
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
