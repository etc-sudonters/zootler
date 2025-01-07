package vm

import (
	"errors"
	"fmt"
	"sudonters/zootler/magicbeanvm/code"
	"sudonters/zootler/magicbeanvm/compiler"
	"sudonters/zootler/magicbeanvm/objects"
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
	stack   stack[objects.Object]
}

func New(objs *objects.Table) VM {
	var r VM
	r.objects = objs
	r.stack = newstack[objects.Object](256)
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

func (this *VM) Execute(executing ExecutionUnit) (result objects.Object, err error) {
	this.stack.reset()
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
		thisOp := code.Op(executing.readu8())
		switch thisOp {
		case code.NOP:
			continue
		case code.ERR:
			err = errors.Join(errors.New("execution halted"), err)
			break loop
		case code.PUSH_T:
			this.stack.push(objects.Boolean(true))
		case code.PUSH_F:
			this.stack.push(objects.Boolean(false))
		case code.PUSH_CONST:
			index := executing.readu16()
			constant := this.objects.Constant(objects.Index(index))
			this.stack.push(constant)
		case code.PUSH_TOKEN:
			ptr := executing.readu16()
			this.stack.push(objects.Pointer(ptr, objects.PtrToken))
		case code.PUSH_SETTING:
			ptr := executing.readu16()
			this.stack.push(objects.Pointer(ptr, objects.PtrSetting))
		case code.PUSH_BUILTIN:
			index := executing.readu16()
			builtin := this.objects.BuiltIn(objects.Index(index))
			this.stack.push(builtin)
		case code.INVERT:
			obj := this.stack.pop()
			this.stack.push(objects.Boolean(!this.truthy(obj)))
		case code.NEED_ALL:
			count := int(executing.readu16())
			var reduction objects.Boolean = true
			for _, obj := range this.stackargs(count) {
				if !this.truthy(obj) {
					reduction = false
					break
				}
			}
			this.stack.popN(count)
			this.stack.push(reduction)
		case code.NEED_ANY:
			count := int(executing.readu16())
			var reduction objects.Boolean = false
			for _, obj := range this.stackargs(count) {
				if this.truthy(obj) {
					reduction = true
					break
				}
			}
			this.stack.ptr -= count
			this.stack.push(reduction)
		case code.CHK_QTY:
			var answer objects.Object
			ptr := executing.readu16()
			qty := executing.readu8()
			answer, err = this.objects.BuiltIns.Has([]objects.Object{
				objects.Pointer(ptr, objects.PtrToken), objects.Number(qty),
			})
			if err != nil {
				err = fmt.Errorf("has 0x%04x 0x%02x: %w", ptr, qty, err)
				break loop
			}
			this.stack.push(answer)
		case code.CHK_ALL:
			var answer objects.Object
			count := int(executing.readu16())
			args := this.stackargs(count)
			answer, err = this.objects.BuiltIns.HasEvery(args)
			if err != nil {
				err = fmt.Errorf("has_every 0x%04x: %w", count, err)
				break loop
			}
			this.stack.popN(count)
			if answer != nil {
				this.stack.push(answer)
			}
		case code.CHK_ANY:
			var answer objects.Object
			count := int(executing.readu16())
			args := this.stackargs(count)
			answer, err = this.objects.BuiltIns.HasAnyOf(args)
			if err != nil {
				err = fmt.Errorf("has_anyof 0x%04x: %w", count, err)
				break loop
			}
			this.stack.popN(count)
			if answer != nil {
				this.stack.push(answer)
			}
		case code.IS_CHILD:
			var answer objects.Object
			answer, err = this.objects.BuiltIns.IsChild(nil)
			if err != nil {
				break loop
			}
			this.stack.push(answer)
		case code.IS_ADULT:
			var answer objects.Object
			answer, err = this.objects.BuiltIns.IsAdult(nil)
			if err != nil {
				break loop
			}
			this.stack.push(answer)
		case code.INVOKE:
			fn := this.stack.pop()
			var answer objects.Object
			switch fn := fn.(type) {
			case *objects.BuiltInFunc:
				argCount := int(executing.readu16())
				if fn.Params > -1 && argCount != fn.Params {
					err = fmt.Errorf("%q expects %d arguments, received %d", fn.Name, fn.Params, argCount)
					break loop
				}
				args := this.stackargs(argCount)
				answer, err = fn.Func(args)
				if err != nil {
					err = fmt.Errorf("%q: %w", fn.Name, err)
					break loop
				}
				this.stack.popN(argCount)
			default:
				err = fmt.Errorf("cannot call %s", fn)
				break loop
			}
			if answer != nil {
				this.stack.push(answer)
			}
		case code.CMP_EQ, code.CMP_NQ, code.CMP_LT:
			err = fmt.Errorf("runtime comparison not implemented")
			break loop
		default:
			err = fmt.Errorf("unrecognized op: 0x%02x", thisOp)
			break loop
		}
	}

	if err == nil && this.stack.ptr > 0 {
		result = this.stack.pop()
	}
	return
}

func (this *VM) stackargs(count int) []objects.Object {
	return this.stack.slice(this.stack.ptr-count, count)
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
