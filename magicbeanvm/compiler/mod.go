package compiler

import (
	"fmt"
	"slices"
	"sudonters/zootler/magicbeanvm/ast"
	"sudonters/zootler/magicbeanvm/code"
	"sudonters/zootler/magicbeanvm/objects"
	"sudonters/zootler/magicbeanvm/symbols"
)

type ByteCode struct {
	Tape code.Instructions
}

func (this *ByteCode) concat(tape code.Instructions) int {
	written := len(tape)
	this.Tape = slices.Concat(this.Tape, tape)
	return written
}

func Compile(nodes ast.Node, symbols *symbols.Table, objects *objects.TableBuilder) (ByteCode, error) {
	var compiler compiler
	var bytecode ByteCode
	compiler.symbols = symbols
	compiler.objects = objects
	compiler.code = &bytecode

	visiting := &compiler
	visitor := ast.Visitor{
		AnyOf:      visiting.AnyOf,
		Boolean:    visiting.Boolean,
		Compare:    visiting.Compare,
		Every:      visiting.Every,
		Identifier: visiting.Identifier,
		Invert:     visiting.Invert,
		Invoke:     visiting.Invoke,
		Number:     visiting.Number,
		String:     visiting.String,
	}
	err := visitor.Visit(nodes)
	return bytecode, err
}

type compiler struct {
	tapePtr int
	symbols *symbols.Table
	objects *objects.TableBuilder
	code    *ByteCode
}

func (this *compiler) emit(op code.Op, operands ...int) int {
	startOfInstruction := this.tapePtr
	this.tapePtr += this.code.concat(code.Make(op, operands...))
	return startOfInstruction
}

func (this *compiler) AnyOf(node ast.AnyOf, visit ast.Visiting) error {
	err := visit.All(node)
	if err != nil {
		return err
	}
	this.emit(code.BEAN_NEED_ANY, len(node))
	return nil
}

func (this *compiler) Boolean(node ast.Boolean, visit ast.Visiting) error {
	if node {
		this.emit(code.BEAN_PUSH_T)
	} else {
		this.emit(code.BEAN_PUSH_F)
	}
	return nil
}

func (this *compiler) Compare(node ast.Compare, visit ast.Visiting) error {
	if err := visit.All([]ast.Node{node.RHS, node.LHS}); err != nil {
		return err
	}
	switch node.Op {
	case ast.CompareEq:
		this.emit(code.BEAN_CMP_EQ)
	case ast.CompareNq:
		this.emit(code.BEAN_CMP_NQ)
	case ast.CompareLt:
		this.emit(code.BEAN_CMP_LT)
	default:
		return fmt.Errorf("uncompilable comparison op: %v", node.Op)
	}

	return nil
}

func (this *compiler) Every(node ast.Every, visit ast.Visiting) error {
	err := visit.All(node)
	if err != nil {
		return err
	}
	this.emit(code.BEAN_NEED_ALL, len(node))
	return nil
}

func (this *compiler) Identifier(node ast.Identifier, visit ast.Visiting) error {
	symbol := this.symbols.LookUpByIndex(node.AsIndex())
	switch symbol.Kind {
	case symbols.BUILT_IN:
		index := this.objects.BuiltIn(symbol.Name)
		this.emit(code.BEAN_PUSH_BUILTIN, int(index))
		return nil
	case symbols.TOKEN, symbols.EVENT:
		index := this.objects.Name(symbol.Name)
		this.emit(code.BEAN_PUSH_PTR, int(index))
	default:
		return fmt.Errorf("uncompilable identifier: %s", symbol)
	}
	return nil
}

func (this *compiler) Invert(node ast.Invert, visit ast.Visiting) error {
	if err := visit(node.Inner); err != nil {
		return err
	}
	this.emit(code.BEAN_PUSH_OPP)
	return nil
}

func (this *compiler) Invoke(node ast.Invoke, visit ast.Visiting) error {
	target := ast.LookUpNodeInTable(this.symbols, node.Target)
	if target != nil {
		var fast bool
		switch target.Name {
		case "has":
			what := ast.LookUpNodeInTable(this.symbols, node.Args[0])
			qty, isQty := node.Args[1].(ast.Number)

			if what != nil && isQty {
				fast = true
				ptr := this.objects.Name(what.Name)
				this.emit(code.BEAN_CHK_QTY, int(ptr), int(qty))
			}
		case "has_anyof":
			fast = true
			if argsErr := visit.All(node.Args); argsErr != nil {
				return argsErr
			}
			this.emit(code.BEAN_CHK_ANY, len(node.Args))
		case "has_every":
			fast = true
			if argsErr := visit.All(node.Args); argsErr != nil {
				return argsErr
			}
			this.emit(code.BEAN_CHK_ALL, len(node.Args))
		case "is_adult":
			fast = true
			this.emit(code.BEAN_IS_ADULT)
		case "is_child":
			fast = true
			this.emit(code.BEAN_IS_CHILD)
		}

		if fast {
			return nil
		}
	}

	if argsErr := visit.All(node.Args); argsErr != nil {
		return argsErr
	}

	if targetErr := visit(node.Target); targetErr != nil {
		return targetErr
	}
	this.emit(code.BEAN_CALL, len(node.Args))
	return nil
}

func (this *compiler) Number(node ast.Number, visit ast.Visiting) error {
	obj := objects.Number(node)
	idx := this.objects.Constant(obj)
	this.emit(code.BEAN_PUSH_CONST, int(idx))
	return nil
}

func (this *compiler) String(node ast.String, visit ast.Visiting) error {
	obj := objects.String(node)
	idx := this.objects.Constant(obj)
	this.emit(code.BEAN_PUSH_CONST, int(idx))
	return nil
}
