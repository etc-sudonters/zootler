package vm

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sudonters/zootler/mido/code"
	"sudonters/zootler/mido/compiler"
	"sudonters/zootler/mido/objects"
)

type VM struct {
	Objects *objects.Table
	Funcs   map[objects.Object]objects.BuiltInFunction
}

func (this *VM) Execute(bytecode compiler.Bytecode) (objects.Object, error) {
	var err error
	result := objects.Null
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
			unit.stack.push(objects.True)
		case code.PUSH_F:
			unit.stack.push(objects.False)
		case code.PUSH_CONST:
			index := unit.readIndex()
			unit.stack.push(this.Objects.AtIndex(index))
		case code.INVERT:
			obj := unit.stack.pop()
			unit.stack.push(objects.PackBool((!this.truthy(obj))))
		case code.NEED_ALL:
			count := int(unit.readu16())
			reduction := true
			stackargs := unit.stackargs(count)
			for _, obj := range stackargs {
				if !this.truthy(obj) {
					reduction = false
					break
				}
			}
			unit.stack.popN(count)
			unit.stack.push(objects.PackBool(reduction))
		case code.NEED_ANY:
			count := int(unit.readu16())
			reduction := false
			stackargs := unit.stackargs(count)
			for _, obj := range stackargs {
				if this.truthy(obj) {
					reduction = true
					break
				}
			}
			unit.stack.popN(count)
			unit.stack.push(objects.PackBool(reduction))
		case code.CHK_QTY:
			ptr := unit.readIndex()
			qty := unit.readu8()
			_, _ = ptr, qty
			unit.stack.push(objects.True)
		case code.INVOKE:
			answer := objects.Null
			fn := unit.stack.pop()
			if !fn.Is(objects.Func) {
				err = fmt.Errorf("cannot call %v", fn)
				break loop
			}
			count := int(unit.readu16())
			args := unit.stackargs(count)
			answer, err = this.call(fn, args)
			unit.stack.popN(count)
			if err != nil {
				break loop
			}
			if answer != objects.Null {
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
	return result, err
}

func (this *VM) truthy(_ objects.Object) bool {
	return true
}

func (this *VM) call(callee objects.Object, args []objects.Object) (objects.Object, error) {
	fn, exists := this.Funcs[callee]
	if !exists {
		return objects.Null, fmt.Errorf("%q is not mapped at runtime", callee)
	}
	return fn(this.Objects, args)
}
