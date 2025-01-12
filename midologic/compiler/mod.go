package compiler

import (
	"fmt"
	"slices"
	"sudonters/zootler/midologic/ast"
	"sudonters/zootler/midologic/code"
	"sudonters/zootler/midologic/objects"
	"sudonters/zootler/midologic/symbols"
)

type ByteCode struct {
	Tape code.Instructions
}

func (this *ByteCode) concat(tape code.Instructions) int {
	written := len(tape)
	this.Tape = slices.Concat(this.Tape, tape)
	return written
}

func Compile(nodes ast.Node, symbols *symbols.Table, objects *objects.TableBuilder, fastops FastOps) (ByteCode, error) {
	var compiler compiler
	var bytecode ByteCode
	compiler.symbols = symbols
	compiler.objects = objects
	compiler.code = &bytecode
	compiler.fastops = fastops

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
	fastops FastOps
	code    *ByteCode
}

func (this *compiler) emit(op code.Op, operands ...int) int {
	return this.join(code.Make(op, operands...))
}

func (this *compiler) join(emitted code.Instructions) int {
	startOfInstruction := this.tapePtr
	this.tapePtr += this.code.concat(emitted)
	return startOfInstruction
}

func (this *compiler) AnyOf(node ast.AnyOf, visit ast.Visiting) error {
	err := visit.All(node)
	if err != nil {
		return err
	}
	this.emit(code.NEED_ANY, len(node))
	return nil
}

func (this *compiler) Boolean(node ast.Boolean, visit ast.Visiting) error {
	if node {
		this.emit(code.PUSH_T)
	} else {
		this.emit(code.PUSH_F)
	}
	return nil
}

func (this *compiler) Compare(node ast.Compare, visit ast.Visiting) error {
	if err := visit.All([]ast.Node{node.RHS, node.LHS}); err != nil {
		return err
	}
	switch node.Op {
	case ast.CompareEq:
		this.emit(code.CMP_EQ)
	case ast.CompareNq:
		this.emit(code.CMP_NQ)
	case ast.CompareLt:
		this.emit(code.CMP_LT)
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
	this.emit(code.NEED_ALL, len(node))
	return nil
}

func (this *compiler) Identifier(node ast.Identifier, visit ast.Visiting) error {
	symbol := this.symbols.LookUpByIndex(node.AsIndex())
	switch symbol.Kind {
	case symbols.BUILT_IN_FUNCTION:
		index := this.objects.GetBuiltIn(symbol.Name)
		this.emit(code.PUSH_BUILTIN, int(index))
		return nil
	case symbols.TOKEN, symbols.EVENT:
		index := this.objects.GetPointerFor(symbol.Name)
		this.emit(code.PUSH_TOKEN, int(index))
	case symbols.SETTING:
		index := this.objects.GetPointerFor(symbol.Name)
		this.emit(code.PUSH_SETTING, int(index))
	default:
		return fmt.Errorf("uncompilable identifier: %s", symbol)
	}
	return nil
}

func (this *compiler) Invert(node ast.Invert, visit ast.Visiting) error {
	if err := visit(node.Inner); err != nil {
		return err
	}
	this.emit(code.INVERT)
	return nil
}

func (this *compiler) Invoke(node ast.Invoke, visit ast.Visiting) error {
	if callee := ast.LookUpNodeInTable(this.symbols, node.Target); callee != nil {
		if fastOp := this.fastops[callee.Name]; fastOp != nil {
			code, err := fastOp(node, this.symbols, this.objects, visit)
			if err != nil {
				return fmt.Errorf("during fastop generation %q: %w", callee.Name, err)
			}

			if len(code) != 0 {
				this.join(code)
				return nil
			}
		}
	}

	if argsErr := visit.All(node.Args); argsErr != nil {
		return argsErr
	}

	if targetErr := visit(node.Target); targetErr != nil {
		return targetErr
	}
	this.emit(code.INVOKE, len(node.Args))
	return nil
}

func (this *compiler) Number(node ast.Number, visit ast.Visiting) error {
	obj := objects.Number(node)
	idx := this.objects.AddConstant(obj)
	this.emit(code.PUSH_CONST, int(idx))
	return nil
}

func (this *compiler) String(node ast.String, visit ast.Visiting) error {
	obj := objects.String(node)
	idx := this.objects.AddConstant(obj)
	this.emit(code.PUSH_CONST, int(idx))
	return nil
}
