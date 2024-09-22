package runtime

import (
	"sudonters/zootler/icearrow/compiler"
	"sudonters/zootler/icearrow/nan"

	"github.com/etc-sudonters/substrate/skelly/stack"
	"github.com/etc-sudonters/substrate/slipup"
)

type VM struct{}

type VMState interface {
	HasQty(uint32, uint8) bool
	HasAny(...uint32) bool
	HasAll(...uint32) bool
	HasBottle() bool
	IsAdult() bool
	IsChild() bool
	AtTod(uint8) bool
}

type Execution struct {
	Result bool
	Err    error
}

func (vm *VM) Execute(tape *compiler.Tape, state VMState, st *compiler.SymbolTable) Execution {
	stk := stack.Make[nan.PackedValue](0, 32)
	PC := 0
	end := tape.Len()

	for PC < end {
		op := compiler.IceArrowOp(tape.Ops[PC])
		switch op {
		case compiler.IA_LOAD_CONST:
			lo := tape.Ops[PC+1]
			hi := tape.Ops[PC+2]
			handle := uint32(hi)<<8 | uint32(lo)
			val := st.Const(handle)
			stk.Push(val.Value)
			PC += 2
			break
		case compiler.IA_LOAD_SYMBOL:
			lo := tape.Ops[PC+1]
			hi := tape.Ops[PC+2]
			handle := uint32(hi)<<8 | uint32(lo)
			val := st.Symbol(handle)
			stk.Push(nan.PackPtr(val.Id))
			PC += 2
			break
		case compiler.IA_LOAD_TRUE:
			stk.Push(nan.PackBool(true))
			break
		case compiler.IA_LOAD_FALSE:
			stk.Push(nan.PackBool(false))
			break
		case compiler.IA_LOAD_IMMED:
			stk.Push(nan.PackUint(tape.Ops[PC+1]))
			PC += 1
			break
		case compiler.IA_LOAD_IMMED2:
			lo := tape.Ops[PC+1]
			hi := tape.Ops[PC+2]
			val := uint16(hi)<<8 | uint16(lo)
			stk.Push(nan.PackUint(val))
			PC += 2
			break
		case compiler.IA_REDUCE_ALL:
			howMany, err := stk.Pop()
			vm.abortEmptyStack(err)

			qty, wasUint := howMany.Uint()
			if !wasUint {
				vm.abortTypeAssert("Uint", howMany)
			}
			result := true
			for range qty {
				v, err := stk.Pop()
				vm.abortEmptyStack(err)

				b, wasBool := v.Bool()
				if !wasBool {
					vm.abortTypeAssert("bool", v)
				}
				result = result && b
			}
			stk.Push(nan.PackBool(result))
			break
		case compiler.IA_REDUCE_ANY:
			howMany, err := stk.Pop()
			vm.abortEmptyStack(err)

			qty, wasUint := howMany.Uint()
			if !wasUint {
				vm.abortTypeAssert("Uint", howMany)
			}
			result := true
			for range qty {
				v, err := stk.Pop()
				vm.abortEmptyStack(err)

				b, wasBool := v.Bool()
				if !wasBool {
					vm.abortTypeAssert("bool", v)
				}
				result = result || b
			}
			stk.Push(nan.PackBool(result))
			break
		case compiler.IA_HAS_QTY:
			lo, hi, qty := tape.Ops[PC+1], tape.Ops[PC+2], tape.Ops[PC+3]
			handle := uint32(hi)<<8 | uint32(lo)
			stk.Push(nan.PackBool(state.HasQty(handle, qty)))
			PC += 3
			break
		case compiler.IA_HAS_ALL:
			howMany, err := stk.Pop()
			vm.abortEmptyStack(err)

			qty, wasUint := howMany.Uint()
			if !wasUint {
				vm.abortTypeAssert("Uint", howMany)
			}
			result := true
			for range qty {
				v, err := stk.Pop()
				vm.abortEmptyStack(err)

				handle, wasTok := v.Token()
				if !wasTok {
					vm.abortTypeAssert("token", v)
				}
				result = result && state.HasQty(handle, 1)
			}
			stk.Push(nan.PackBool(result))
			break
		case compiler.IA_HAS_ANY:
			howMany, err := stk.Pop()
			vm.abortEmptyStack(err)

			qty, wasUint := howMany.Uint()
			if !wasUint {
				vm.abortTypeAssert("Uint", howMany)
			}
			result := true
			for range qty {
				v, err := stk.Pop()
				vm.abortEmptyStack(err)

				handle, wasTok := v.Token()
				if !wasTok {
					vm.abortTypeAssert("token", v)
				}
				result = result || state.HasQty(handle, 1)
			}
			stk.Push(nan.PackBool(result))
			break
		case compiler.IA_IS_ADULT:
			stk.Push(nan.PackBool(state.IsAdult()))
			break
		case compiler.IA_IS_CHILD:
			stk.Push(nan.PackBool(state.IsChild()))
			break
		case compiler.IA_HAS_BOTTLE:
			stk.Push(nan.PackBool(state.HasBottle()))
			break
		case compiler.IA_CHK_TOD:
			tod := tape.Ops[PC+1]
			stk.Push(nan.PackBool(state.AtTod(tod)))
			PC += 1
			break
		default:
			panic(slipup.Createf("unknown IAOP: 0x%02X", op))
		}
		PC += 1
	}

	result, err := stk.Pop()
	if err != nil {
		err = slipup.Describe(err, "could not set execution result")
		return Execution{false, err}
	}

	answer, wasBool := result.Bool()
	if !wasBool {
		return Execution{false, slipup.Createf("expected bool result")}
	}

	return Execution{answer, err}
}

func (vm *VM) abortOnErr(err error, tpl string, v ...any) {
	if err != nil {
		panic(slipup.Describef(err, tpl, v...))
	}
}

func (vm *VM) abortEmptyStack(err error) {
	vm.abortOnErr(err, "no values pushed to stack")
}

func (vm *VM) abortTypeAssert(expected string, pv nan.PackedValue) {
	panic(slipup.Createf("expected %q, got %s", expected, pv))
}
