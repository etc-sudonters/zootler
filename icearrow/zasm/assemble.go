package zasm

import (
	"errors"
	"sudonters/zootler/icearrow/ast"

	"github.com/etc-sudonters/substrate/slipup"
)

type Assembler struct {
	Data DataBuilder
}

func (a *Assembler) CreateDataTables() Data {
	return CreateDataTables(a.Data)
}

type Assembly struct {
	units []Unit
	data  Data
}

func (a *Assembly) Include(assembly Unit) {
	a.units = append(a.units, assembly)
}

func (a *Assembly) AttachDataTables(data Data) {
	a.data = data
}

type Unit struct {
	Name string
	I    Instructions
}

func ToZasmOp[u ~uint8](op u) Op {
	switch uint8(op) {
	case uint8(ast.AST_CMP_EQ):
		return OP_CMP_EQ
	case uint8(ast.AST_CMP_NQ):
		return OP_CMP_NQ
	case uint8(ast.AST_CMP_LT):
		return OP_CMP_LT
	case uint8(ast.AST_BOOL_AND):
		return OP_BOOL_AND
	case uint8(ast.AST_BOOL_OR):
		return OP_BOOL_OR
	case uint8(ast.AST_BOOL_NEGATE):
		return OP_BOOL_NEGATE
	default:
		panic(slipup.Createf("unknown ast operation %v", op))
	}
}

func (a *Assembler) Assemble(label string, tree ast.Node) (Unit, error) {
	var ass Unit
	ass.Name = label
	i, err := ast.Transform(a, tree)
	ass.I = i
	return ass, err
}

func (a *Assembler) Comparison(node *ast.Comparison) (Instructions, error) {
	lhs, lhErr := ast.Transform(a, node.LHS)
	rhs, rhErr := ast.Transform(a, node.RHS)
	if joined := errors.Join(lhErr, rhErr); joined != nil {
		panic(joined)
	}
	if lhs == nil || rhs == nil {
		panic(slipup.Createf("expected both arms to exist"))
	}

	return IntoIW(lhs).Union(rhs).WriteOp(ToZasmOp(node.Op)).I, nil
}

func (a *Assembler) BooleanOp(node *ast.BooleanOp) (Instructions, error) {
	if node.Op == ast.AST_BOOL_NEGATE {
		return a.negate(node.LHS)
	}

	lhs, lhErr := ast.Transform(a, node.LHS)
	rhs, rhErr := ast.Transform(a, node.RHS)
	if joined := errors.Join(lhErr, rhErr); joined != nil {
		panic(joined)
	}
	if lhs == nil || rhs == nil {
		panic(slipup.Createf("expected both arms to exist"))
	}
	return IntoIW(lhs).Union(rhs).WriteOp(ToZasmOp(node.Op)).I, nil
}

func (a *Assembler) Call(node *ast.Call) (Instructions, error) {
	iw := IntoIW(a.callArgs(node.Args))
	name := a.Data.Names.Intern(node.Callee)
	return iw.WriteCall(name, len(node.Args)).I, nil
}

func (a *Assembler) Identifier(node *ast.Identifier) (Instructions, error) {
	name := a.Data.Names.Intern(node.Name)
	return IW().WriteLoadIdent(name).I, nil
}

func (a *Assembler) Literal(node *ast.Literal) (Instructions, error) {
	switch v := node.Value.(type) {
	case bool:
		return IW().WriteLoadBool(v).I, nil
	case float64:
		c := a.Data.Consts.Intern(Pack(v))
		return IW().WriteLoadConst(c).I, nil
	case string:
		str := a.Data.Strs.Intern(v)
		return IW().WriteLoadStr(str).I, nil
	default:
		panic("unreachable")
	}
}

func (a *Assembler) Empty(node *ast.Empty) (Instructions, error) {
	return nil, nil
}

func (a *Assembler) negate(node ast.Node) (Instructions, error) {
	instructions, _ := ast.Transform(a, node)
	return IntoIW(instructions).WriteNegate().I, nil
}

func (a *Assembler) callArgs(args []ast.Node) Instructions {
	iw := IW()
	for _, arg := range args {
		instrs, _ := ast.Transform(a, arg)
		iw = iw.Union(instrs)
	}
	return iw.I
}
